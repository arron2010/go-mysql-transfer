package endpoint

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"go-mysql-transfer/model"
	"sync"
)

var (
	_myMQ *myMQEndpoint
	once  sync.Once
)

type myMQEndpoint struct {
}

func newMyMQEndpoint() *myMQEndpoint {
	//m := &myMQEndpoint{}
	once.Do(func() {
		_myMQ = &myMQEndpoint{}
	})
	return _myMQ
}

func (mq *myMQEndpoint) Connect() error {
	return nil
}

func (mq *myMQEndpoint) Ping() error {
	return nil
}

func (mq *myMQEndpoint) Consume(serverName string, position mysql.Position, rows []*model.RowRequest) error {
	for _, row := range rows {
		printRow(serverName, row, position)
	}
	return nil
}

func (mq *myMQEndpoint) Stock(requests []*model.RowRequest) int64 {
	return 0
}

func (mq *myMQEndpoint) Close() {

}

func printRow(serverName string, row *model.RowRequest, position mysql.Position) {
	var rowText string
	rowText = fmt.Sprintf("serverName:%s Action:%s RuleKey:%s  Values:%v Position:%v", serverName, row.Action, row.RuleKey, row.Row, position)
	fmt.Println(rowText)
}
