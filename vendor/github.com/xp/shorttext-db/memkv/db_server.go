package memkv

import (
	"github.com/xp/shorttext-db/config"

	"github.com/xp/shorttext-db/server"
)

type MemDBServer struct {
	node *server.Node

	Id int
	db MemDB
}

func NewDBServer(node *server.Node) *MemDBServer {
	var err error
	server := &MemDBServer{}
	c := config.GetCase()
	id := int(c.Local.ID)
	server.db, err = CreatDB(uint32(id), false)
	server.db.CreateIndex("KeyIndex", "*", IndexKey)
	server.db.CreateIndex("IndexRawKey", "*", IndexRawKey)

	if err != nil {
		panic(err)
	}
	server.db.SetId(uint32(id))
	server.Id = id
	node.RegisterHandler(server)
	server.node = node
	initialize(server.db)
	return server
}

func (s *MemDBServer) Handle(msgType uint32, data []byte) ([]byte, bool, error) {
	var err error
	var resp []byte

	switch msgType {
	case config.MSG_KV_SET:
		dbItem := &DBItem{}
		err = unmarshalDbItem(data, dbItem)
		if err != nil {
			return nil, true, err
		}
		s.debugOp("Put", dbItem)
		dbItem.Key = mvccEncode(dbItem.RawKey, dbItem.CommitTS)
		err = s.db.Put(dbItem)
		//s.db.Ascend("IndexRawKey", func(key Key, value *DBItem) bool {
		//	fmt.Println(value.RawKey)
		//	return true
		//});

		return nil, true, err

	case config.MSG_KV_FIND:
		dbItem := &DBItem{}
		err = unmarshalDbItem(data, dbItem)
		if err != nil {
			return nil, true, err
		}
		//s.debugOp("Find", dbItem)
		items := []*DBItem{} //s.db.Scan(dbItem.Key, dbItem.Value)
		protoItems := &DBItems{}
		protoItems.Items = items
		resp, err = marshalDbItems(protoItems)
		return resp, true, err
	case config.MSG_KV_DEL:
		dbItem := &DBItem{}
		err = unmarshalDbItem(data, dbItem)
		if err != nil {
			return nil, true, err
		}
		err = s.db.Delete(dbItem.Key)
		return nil, true, err
	}

	return resp, false, err
}

func (s *MemDBServer) ReportUnreachable(id uint64) {

}
func (s *MemDBServer) debugOp(op string, dbItem *DBItem) {

	logger.Infof("%s DbItem Key[%v] Ts[%d]\n", op, dbItem.RawKey, dbItem.CommitTS)

}
