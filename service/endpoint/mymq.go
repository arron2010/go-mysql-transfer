package endpoint

import (
	"fmt"
	"github.com/asim/mq/go/client"
	mqgrpc "github.com/asim/mq/go/client/grpc"
	"github.com/asim/mq/go/client/selector"
	proto2 "github.com/golang/protobuf/proto"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"github.com/siddontang/go-mysql/schema"
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
	mqAddrs  []string
	//isFirst bool
	consumePool *gopool.PoolWithFunc
}

func newMyMQEndpoint() *myMQEndpoint {
	once.Do(func() {
		_myMQ = &myMQEndpoint{}
		//_myMQ.isFirst = true
	})
	return _myMQ
}

func (mq *myMQEndpoint) Connect() error {
	config := global.Cfg()
	mq.mqAddrs = strings.Split(config.MQAddr, ",")
	options := gopool.Options{}
	options.ExpiryDuration = time.Duration(3) * time.Second
	options.Nonblocking = false
	options.PreAlloc = true
	consumePool, err := gopool.NewPoolWithFunc(config.MQPoolSize, publish, gopool.WithOptions(options))
	mq.consumePool = consumePool
	if err != nil {
		return err
	}

	mq.mqClient = mqgrpc.New(client.WithServers(mq.mqAddrs...),
		client.WithSelector(&selector.Shard{}),
		client.WithRetries(3))

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
	//for _, row := range rows {
	//	printRow(serverName, row, position)
	//}
	return nil
}

func (mq *myMQEndpoint) Stock(requests []*model.RowRequest) int64 {
	return 0
}

func (mq *myMQEndpoint) Close() {
	mq.mqClient.Close()
}

func publish(reqObj interface{}) {
	var req *mqRequest
	req = reqObj.(*mqRequest)
	l := len(req.rows)
	for i := 0; i < l; i++ {
		topic := global.BuildRootPath(req.serverName, req.rows[i].RuleKey)

		ev := req.rows[i].ReplicationRowsEvent
		buf, err := buildMessage(req.serverName, req.rows[i], req.ruleMap, ev)
		if err != nil {
			logs.Errorf("数据传递时，序列化错误:%s %s %v", topic, req.rows[i].Action, req.rows[i].Raw)
			break
		}
		//logs.Infof("消息大小:%d",len(buf))

		//r,_ :=ev.DecodeRows(req.rows[i].Raw,ev.Table,ev.ColumnBitmap1)
		//fmt.Println(r)
		err = req.mq.mqClient.Publish(uint64(req.rows[i].Timestamp), topic, buf)
		if err != nil {
			logs.Errorf("数据传递时，订阅消息错误:%s %s %v", topic, req.rows[i].Action, req.rows[i].Row)
		}
		rowText := fmt.Sprintf("Position:%s Server:%s Action:%s RuleKey:%s  Row Data Size:%v",
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
	rule, ok := ruleMap.RuleIns(row.RuleKey)
	if ok {
		tbl := rule.TableInfo
		buildPKColumns(bufRow, tbl.PKColumns)
		buildColumns(bufRow, tbl.Columns)
	}
	if len(row.Raw) > 0 {
		buildColumnValue(bufRow, row)
	}
	bufRow.Val = row.Raw
	bufRow.ColumnMeta = collections.Uint16ToUInt32(ev.Table.ColumnMeta)
	bufRow.ColumnType = ev.Table.ColumnType
	bufRow.NeedBitmap2 = boolutil.BoolToInt(ev.NeedBitmap2)
	bufRow.ColumnBitmap1 = ev.ColumnBitmap1
	bufRow.ColumnBitmap2 = ev.ColumnBitmap2

	//r,_ := encoding.DecodeRows(row.Raw,ev.ColumnCount,ev.ColumnBitmap1,ev.ColumnBitmap2,
	//	ev.Table.ColumnType,bufRow.ColumnMeta,ev.NeedBitmap2)
	//fmt.Println(r)

	//if len(row.Old) > 0 {
	//	bufRow.OldValues = buildColumnValue(bufRow, row.Old)
	//}
	buf, err = proto2.Marshal(bufRow)
	return buf, err
}

func buildPKColumns(bufRow *proto.Row, cols []int) {
	bufRow.PKColumns = make([]uint32, 0, len(cols))
	for i := 0; i < len(cols); i++ {
		bufRow.PKColumns = append(bufRow.PKColumns, uint32(cols[i]))
	}
}

func buildColumns(bufRow *proto.Row, cols []schema.TableColumn) {
	bufRow.Columns = make([]*proto.ColumnInfo, 0, len(cols))
	for i := 0; i < len(cols); i++ {
		bufRow.Columns = append(bufRow.Columns, &proto.ColumnInfo{
			Name:       cols[i].Name,
			Type:       uint32(cols[i].Type),
			Collation:  cols[i].Collation,
			RawType:    cols[i].RawType,
			IsAuto:     boolutil.BoolToInt(cols[i].IsAuto),
			IsUnsigned: boolutil.BoolToInt(cols[i].IsUnsigned),
			IsVirtual:  boolutil.BoolToInt(cols[i].IsVirtual),
			FixedSize:  uint32(cols[i].FixedSize),
			MaxSize:    uint32(cols[i].MaxSize),
		})
	}
}

func buildColumnValue(bufRow *proto.Row, row *model.RowRequest) {
	bufRow.Val = row.Raw

	//var colValue *proto.ColumnValue
	//if len(bufRow.Columns) != len(values) {
	//	return nil
	//}
	//columnValues := make([]*proto.ColumnValue, 0, len(values))
	//for i := 0; i < len(values); i++ {
	//	colValue = createColumnValue(bufRow.Columns[i], values[i])
	//	columnValues = append(columnValues, colValue)
	//}
	//return columnValues
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
