package endpoint

import (
	"errors"
	"fmt"
	"github.com/asim/mq/common"
	"github.com/asim/mq/go/client"
	proto2 "github.com/golang/protobuf/proto"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"github.com/siddontang/go-mysql/schema"
	"github.com/xp/shorttext-db/network/proxy"
	"go-mysql-transfer/global"
	"go-mysql-transfer/gopool"
	"go-mysql-transfer/model"
	"go-mysql-transfer/proto"
	"go-mysql-transfer/util/boolutil"
	"go-mysql-transfer/util/byteutil"
	"go-mysql-transfer/util/collections"
	"go-mysql-transfer/util/logs"
	"strings"
	"sync"
	"time"
)

var (
	_myMQ *myMQEndpoint
	once  sync.Once
)

type mqRequest struct {
	serverName string
	rows       []*model.RowRequest
	ruleMap    *global.RuleMap
	mq         *myMQEndpoint
	position   mysql.Position
}

type myMQEndpoint struct {
	mqClient client.Client
	nodeIds  []uint64
	//isFirst bool
	node        *proxy.NodeProxy
	consumePool *gopool.PoolWithFunc
}

func newMyMQEndpoint() *myMQEndpoint {
	once.Do(func() {
		_myMQ = &myMQEndpoint{}
		cfg := global.Cfg()
		peers := strings.Split(cfg.Peers, ",")
		_myMQ.nodeIds = getNodeIds(len(peers))
		_myMQ.node = proxy.NewNodeProxy(peers, cfg.LogLevel)
		//_myMQ.isFirst = true
	})
	return _myMQ
}
func getNodeIds(count int) []uint64 {
	ids := make([]uint64, 0, count)
	for i := 2; i <= count; i++ {
		ids = append(ids, uint64(i))
	}
	return ids
}
func (mq *myMQEndpoint) Connect() error {
	config := global.Cfg()

	options := gopool.Options{}
	options.ExpiryDuration = time.Duration(3) * time.Second
	options.Nonblocking = false
	options.PreAlloc = true
	consumePool, err := gopool.NewPoolWithFunc(config.MQPoolSize, publish, gopool.WithOptions(options))
	mq.consumePool = consumePool
	if err != nil {
		return err
	}

	//mq.mqClient = mqgrpc.New(client.WithServers(mq.mqAddrs...),
	//	client.WithSelector(&selector.Shard{}),
	//	client.WithRetries(3))

	return nil
}

func (mq *myMQEndpoint) Ping() error {
	_, err := mq.mqClient.Ping("mq")
	return err
}

func (mq *myMQEndpoint) Consume(serverName string, position mysql.Position, rows []*model.RowRequest, ruleMap *global.RuleMap) error {
	req := &mqRequest{serverName: serverName, rows: rows, ruleMap: ruleMap, mq: mq, position: position}
	err := mq.consumePool.Invoke(req)
	if err != nil {
		return err
	}
	return nil
}

func (mq *myMQEndpoint) Stock(requests []*model.RowRequest) int64 {
	return 0
}

func (mq *myMQEndpoint) Close() {
	mq.mqClient.Close()
}

const _sendRetries = 3

func (mq *myMQEndpoint) pub(timestamp uint64, topic string, data []byte) error {
	index := timestamp % uint64(len(mq.nodeIds))
	var to uint64
	for i := 0; i < _sendRetries; i++ {
		to = mq.nodeIds[int(index)]
		alive := mq.node.IsAlive(to)
		if !alive {
			continue
		}
		_, err := mq.node.SendKeyMsg(topic, to, common.OP_PUB, data)
		return err
	}
	return errors.New(fmt.Sprintf("节点【%d 】连接中断", to))
}
func publish(reqObj interface{}) {
	var req *mqRequest
	req = reqObj.(*mqRequest)
	l := len(req.rows)
	for i := 0; i < l; i++ {
		ev := req.rows[i].ReplicationRowsEvent

		topic := global.BuildRootPath(req.serverName, req.rows[i].RuleKey)
		buf, err := buildMessage(req.serverName, req.rows[i], req.ruleMap, ev)
		if err != nil {
			logs.Errorf("数据传递时，序列化错误:%s %s %v", topic, req.rows[i].Action, req.rows[i].Raw)
			break
		}

		err = req.mq.pub(uint64(req.rows[i].Timestamp), topic, buf)
		if err != nil {
			logs.Errorf("数据传递时，订阅消息错误:%s %s %v", topic, req.rows[i].Action, err)
		}
		rowText := fmt.Sprintf("Position:%s Server:%s Action:%s RuleKey:%s  DataSize:%v",
			req.position.Name, req.serverName, req.rows[i].Action, topic, len(buf))
		logs.Infof("消息发送>>%s", rowText)
	}
}

func buildMessage(serverName string, row *model.RowRequest, ruleMap *global.RuleMap, ev *replication.RowsEvent) ([]byte, error) {
	var bufRow *proto.Row
	var buf []byte
	var err error
	bufRow = &proto.Row{}
	dbInfo := strings.Split(row.RuleKey, "/")
	bufRow.Server = serverName
	bufRow.Action = row.Action
	bufRow.DB = dbInfo[0]
	bufRow.Table = dbInfo[1]
	bufRow.Timestamp = uint64(row.Timestamp)

	bufRow.Val = row.Raw
	bufRow.ColumnMeta = collections.Uint16ToUInt32(ev.Table.ColumnMeta)
	bufRow.ColumnType = ev.Table.ColumnType
	bufRow.NeedBitmap2 = boolutil.BoolToInt(ev.NeedBitmap2)
	bufRow.ColumnBitmap1 = ev.ColumnBitmap1
	bufRow.ColumnBitmap2 = ev.ColumnBitmap2
	bufRow.ColumnCount = uint32(ev.ColumnCount)

	//needBitmap2 := false
	//if bufRow.NeedBitmap2 > 0 {
	//	needBitmap2 = true
	//}
	//r, _ := encoding.DecodeRows(bufRow.Val, uint64(len(bufRow.Columns)), bufRow.ColumnBitmap1, bufRow.ColumnBitmap2,
	//	bufRow.ColumnType, bufRow.ColumnMeta, needBitmap2)
	//
	//fmt.Println(r)

	buf, err = proto2.Marshal(bufRow)
	return buf, err
}

func buildPKColumns(bufRow *proto.Row, cols []int) {
	bufRow.PKColumns = make([]uint32, 0, len(cols))
	for i := 0; i < len(cols); i++ {
		bufRow.PKColumns = append(bufRow.PKColumns, uint32(cols[i]))
	}
}

func buildColumnValue(bufRow *proto.Row, row *model.RowRequest) {
	bufRow.Val = row.Raw
}

func createColumnValue(col *proto.ColumnInfo, value interface{}) *proto.ColumnValue {
	colValue := &proto.ColumnValue{}
	if value == nil {
		return colValue
	}
	var buf []byte
	switch col.Type {
	case schema.TYPE_ENUM, schema.TYPE_SET, schema.TYPE_NUMBER, schema.TYPE_DECIMAL, schema.TYPE_FLOAT:
		switch val := value.(type) {
		case int64:
			buf = byteutil.Int64ToBytes(val)
		case string:
			buf = []byte(val)
		case []byte:
			buf = val
		default:
			buf = nil
		}
	case schema.TYPE_BIT:
		switch val := value.(type) {
		case string:
			if val == "\x01" {
				buf = byteutil.Int64ToBytes(1)
			}
			buf = byteutil.Int64ToBytes(0)
		default:
			buf = nil
		}
	case schema.TYPE_STRING, schema.TYPE_JSON, schema.TYPE_DATETIME, schema.TYPE_TIMESTAMP, schema.TYPE_DATE:
		switch val := value.(type) {
		case string:
			buf = []byte(val)
		case []byte:
			buf = val
		default:
			buf = nil
		}
	default:
		buf = nil
	}
	colValue.Val = buf
	return colValue
}

func printRow(serverName string, row *model.RowRequest, position mysql.Position) {
	var rowText string
	rowText = fmt.Sprintf("serverName:%s Action:%s RuleKey:%s  Values:%v Position:%v", serverName, row.Action, row.RuleKey, row.Row, position)
	fmt.Println(rowText)
}
