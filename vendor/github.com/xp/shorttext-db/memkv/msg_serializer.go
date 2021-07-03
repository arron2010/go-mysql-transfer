package memkv

import (
	proto2 "github.com/golang/protobuf/proto"
)

type MsgSerializer struct {
}

func NewMsgSerializer() *MsgSerializer {
	m := &MsgSerializer{}
	return m
}
func (m *MsgSerializer) Serialize(source interface{}) ([]byte, error) {
	var buf []byte
	var err error
	switch source.(type) {
	case *DBQueryParam:
		obj := source.(*DBQueryParam)
		buf, err = proto2.Marshal(obj)
	case *DBItems:
		obj := source.(*DBItems)
		buf, err = proto2.Marshal(obj)
	}
	return buf, err
}

func (m *MsgSerializer) Deserialize(typeName string, payload []byte) (interface{}, error) {
	var err error
	switch typeName {
	case "*memkv.DBQueryParam":
		obj := &DBQueryParam{}
		err = proto2.Unmarshal(payload, obj)
		return obj, err
	case "*memkv.DBItems":
		obj := &DBItems{}
		err = proto2.Unmarshal(payload, obj)
		return obj, err
	}
	return nil, nil
}
