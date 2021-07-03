package memkv

import (
	"github.com/xp/shorttext-db/config"
	"github.com/xp/shorttext-db/easymr/artifacts/task"
	"github.com/xp/shorttext-db/easymr/collaborator"
	"github.com/xp/shorttext-db/easymr/interfaces"
	"github.com/xp/shorttext-db/network/proxy"
	"github.com/xp/shorttext-db/server"
	"sync"
)

type RemoteDBProxy struct {
	clbt *collaborator.Collaborator
	n    *proxy.NodeProxy
	c    *chooser
}

var remoteOnce sync.Once
var remoteDBProxy *RemoteDBProxy

func GetRemoteDBProxy() *RemoteDBProxy {
	remoteOnce.Do(func() {
		remoteDBProxy = NewRemoteDBProxy(server.GetNodeProxy(), collaborator.GetCollaborator())
	})
	return remoteDBProxy
}
func NewRemoteDBProxy(n *proxy.NodeProxy, clbt *collaborator.Collaborator) *RemoteDBProxy {
	r := &RemoteDBProxy{}
	c := config.GetCase()
	r.n = n
	r.c = NewChooser()
	r.c.masterId = uint32(n.Id)
	r.clbt = clbt
	r.c.SetBuckets(c.GetCardList())
	r.c.n = n
	initialize(nil)

	return r
}
func (r *RemoteDBProxy) FindByKey(finding Key, locked bool) []*DBItem {
	return nil
}

func (r *RemoteDBProxy) Scan(startKey Key, endKey Key, ts uint64, limit int, desc bool, validate ValidateFunc) []*DBItem {
	var nDesc uint32
	if desc {
		nDesc = DESCEND_RANGE
	} else {
		nDesc = ASCEND_RANGE
	}
	jobInfo := interfaces.NewSimpleJobInfo("MemKVJob", false,
		&DBQueryParam{StartKey: startKey, EndKey: endKey, Ts: ts, Limit: uint32(limit), Desc: nDesc})
	context := &task.TaskContext{}
	context.Context = make(map[string]interface{})
	result, err := r.clbt.MapReduce(jobInfo, context)
	if err != nil {
		logger.Errorf("RemoteDBProxy区间查询错误:%v startKey:%v endKey:%v ts:%d limit:%d desc:%v\n",
			err, startKey, endKey, ts, limit, desc)
		return nil
	} else {
		dbItems := result.Content.(*DBItems)
		return dbItems.Items
	}
}

func (r *RemoteDBProxy) GetByRawKey(key []byte, ts uint64) (*DBItem, bool) {
	return r.findValue(key, ts, GET_BY_RAW_KEY)
}

func (r *RemoteDBProxy) Put(item *DBItem) (err error) {
	var to uint64

	to, _ = r.c.Choose(item.RawKey, true)
	_, err = r.send(item, to, config.MSG_KV_SET)
	if err == nil {
		r.c.UpdateRegion(to, 0, 1)
	}
	return err
}
func (r *RemoteDBProxy) findValue(key []byte, ts uint64, op uint32) (item *DBItem, validated bool) {
	jobInfo := interfaces.NewSimpleJobInfo("MemKVJob", false,
		&DBQueryParam{StartKey: key, EndKey: nil, Ts: ts, Limit: 1, Desc: op})
	context := &task.TaskContext{}
	context.Context = make(map[string]interface{})
	result, err := r.clbt.MapReduce(jobInfo, context)

	if err != nil {
		logger.Errorf("RemoteDBProxy按键查询错误:%v key:%v ts:%d\n",
			err, key, ts)
		return nil, false
	} else {
		if result == nil {
			return nil, false
		}
		dbItems := result.Content.(*DBItems)
		if len(dbItems.Items) > 0 {
			return dbItems.Items[0], len(dbItems.Items[0].RawKey) > 0
		}
		return nil, validated
	}
}
func (r *RemoteDBProxy) Get(key []byte, ts uint64) (item *DBItem, validated bool) {

	return r.findValue(key, ts, GET)
}

func (r *RemoteDBProxy) Delete(key []byte, ts uint64) (err error) {
	jobInfo := interfaces.NewSimpleJobInfo("MemKVJob", false,
		&DBQueryParam{StartKey: key, EndKey: nil, Ts: ts, Limit: 1, Desc: DELETE})
	context := &task.TaskContext{}
	context.Context = make(map[string]interface{})
	_, err = r.clbt.MapReduce(jobInfo, context)
	return err
}

func (r *RemoteDBProxy) Close() error {
	return r.c.Close()
}

func (r *RemoteDBProxy) find(item *DBItem) (result *DBItems, err error) {
	to, _ := r.c.Choose(item.Key, false)
	if len(item.Key) > 0 {
		item.Key = mvccEncode(item.Key, lockVer)
		item.Val = mvccEncode(item.Val, 0)
	}
	if to == 0 {
		return NewDbItems(), nil
	}

	result, err = r.send(item, to, config.MSG_KV_FIND)
	if result == nil {
		result = NewDbItems()
	}
	return result, err
}

func (r *RemoteDBProxy) send(item *DBItem, to uint64, op uint32) (items *DBItems, err error) {
	var req, resp []byte
	req, err = marshalDbItem(item)
	if err != nil {
		return nil, err
	}

	resp, err = r.n.SendSingleMsg(to, op, req)
	if op == config.MSG_KV_FIND {
		items = NewDbItems()
		err = unmarshalDbItems(resp, items)
		return items, err
	} else {
		return nil, err
	}
}
