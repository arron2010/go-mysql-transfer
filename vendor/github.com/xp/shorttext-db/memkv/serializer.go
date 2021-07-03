package memkv

import (
	proto2 "github.com/golang/protobuf/proto"
)

func marshalDbItem(dbItem *DBItem) ([]byte, error) {
	return proto2.Marshal(dbItem)
}
func unmarshalDbItem(buff []byte, dbItem *DBItem) error {
	return proto2.Unmarshal(buff, dbItem)
}

func unmarshalDbItems(buff []byte, dbItem *DBItems) error {
	return proto2.Unmarshal(buff, dbItem)
}
func marshalDbItems(dbItems *DBItems) ([]byte, error) {
	return proto2.Marshal(dbItems)
}
