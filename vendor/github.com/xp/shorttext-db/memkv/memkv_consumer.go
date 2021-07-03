package memkv

import (
	"bytes"
	"github.com/xp/shorttext-db/easymr/artifacts/task"
	"github.com/xp/shorttext-db/utils"
)

type MemKVConsumer struct {
	db MemDB
}

func NewMemKVConsumer(db MemDB) *MemKVConsumer {
	m := &MemKVConsumer{}
	m.db = db
	return m
}

func (m *MemKVConsumer) handle(queryParam *DBQueryParam) []*DBItem {
	var data []*DBItem
	switch queryParam.Desc {
	case DESCEND_RANGE:
		data = scan(m.db, queryParam.StartKey, queryParam.EndKey, queryParam.Ts, int(queryParam.Limit), true, nil)
	case ASCEND_RANGE:
		data = scan(m.db, queryParam.StartKey, queryParam.EndKey, queryParam.Ts, int(queryParam.Limit), false, nil)
	case GET:
		k := mvccEncode(queryParam.StartKey, queryParam.Ts)
		item := m.db.Get(k)
		data = []*DBItem{item}
	case GET_BY_RAW_KEY:
		pivot := &DBItem{RawKey: queryParam.StartKey, CommitTS: queryParam.Ts}
		m.db.AscendGreaterOrEqual("IndexRawKey", pivot, func(key Key, value *DBItem) bool {
			if pivot.CommitTS >= value.CommitTS && bytes.Compare(pivot.RawKey, value.RawKey) == 0 {
				data = []*DBItem{value}
				return false
			}
			return true
		})
	case DELETE:
		k := mvccEncode(queryParam.StartKey, queryParam.Ts)
		m.db.Delete(k)
	}
	return data
}
func (m *MemKVConsumer) Consume(workerId uint, taskItem *task.Task) bool {
	logger.Info("开始本地处理")
	if utils.IsNil(m.db) {
		taskItem.Object = nil
		return true
	}
	if taskItem.Object == nil {
		return true
	}
	queryParam := taskItem.Object.(*DBQueryParam)
	//desc := false
	//if queryParam.Desc == DESCEND_RANGE{
	//	desc= true
	//}
	//data :=scan(m.db,queryParam.StartKey, queryParam.EndKey,queryParam.Ts,int(queryParam.Limit),desc,nil)
	data := m.handle(queryParam)

	if len(data) > 0 {
		result := &DBItems{Items: data}
		taskItem.Result.Append(result)
	}

	//TODO 后续转换成PROTO对象
	logger.Infof("完成本地处理,记录数:%d\n", len(data))

	return true
}
