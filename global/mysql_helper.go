package global

import (
	"github.com/asim/mq/glogger"
	"github.com/pingcap/errors"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
)

type MySQLHelper struct {
	cfg  *ServerConfig
	conn *client.Conn
}

func NewMySQLHelper(config *ServerConfig) *MySQLHelper {
	helper := &MySQLHelper{}
	helper.cfg = config
	return helper
}

func (This *MySQLHelper) execute(cmd string, args ...interface{}) (rr *mysql.Result, err error) {
	retryNum := 3
	for i := 0; i < retryNum; i++ {
		if This.conn == nil {
			This.conn, err = client.Connect(This.cfg.Addr, This.cfg.User, This.cfg.Password, "")
			if err != nil {
				return nil, errors.Trace(err)
			}
		}

		rr, err = This.conn.Execute(cmd, args...)
		if err != nil && !mysql.ErrorEqual(err, mysql.ErrBadConn) {
			return
		} else if mysql.ErrorEqual(err, mysql.ErrBadConn) {
			This.conn.Close()
			This.conn = nil
			continue
		} else {
			return
		}
	}
	return
}

const _sqlGetLastBinlog = `show master status`

func (This *MySQLHelper) GetLastBinlog() (string, uint32, error) {

	resultSet, err := This.execute(_sqlGetLastBinlog)
	if err != nil {
		glogger.Errorf("获取日志文件名错误: %s - %s", _sqlGetLastBinlog, err.Error())
		return "", 0, err
	}
	rowNumber := resultSet.RowNumber()
	val1, err := resultSet.GetString(rowNumber-1, 0)
	val2, err := resultSet.GetUint(rowNumber-1, 1)
	if err != nil {
		glogger.Errorf("数据导出错误: %s - %s", _sqlGetLastBinlog, err.Error())
		return "", 0, err
	}
	return val1, uint32(val2), nil
}
