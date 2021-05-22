package storage

import (
	"github.com/siddontang/go-mysql/mysql"
	"go-mysql-transfer/global"
)

type PositionStorageEx interface {
	Initialize() error
	Save(name string, pos mysql.Position) error
	Get(name string) (mysql.Position, error)
}

func NewPositionStorageEx() PositionStorageEx {
	if global.Cfg().IsCluster() {
		if global.Cfg().IsEtcd() {
			return &etcdPositionStorageEx{}
		}
	}
	return nil
}
