package storage

import (
	"encoding/json"
	"github.com/siddontang/go-mysql/mysql"

	"go-mysql-transfer/global"
	"go-mysql-transfer/util/etcds"
)

type etcdPositionStorageEx struct {
}

func (s *etcdPositionStorageEx) appendName(name string) string {
	return global.Cfg().ZkPositionDir() + "/" + name
}
func (s *etcdPositionStorageEx) Initialize() error {

	for i := 0; i < len(global.Cfg().ServerConfigs); i++ {
		serverConfig := global.Cfg().ServerConfigs[i]
		//positionInfo := strings.Split(serverConfig.Position,".")
		data, err := json.Marshal(mysql.Position{
			Name: serverConfig.Position,
			Pos:  1,
		})
		if err != nil {
			return err
		}

		err = etcds.CreateIfNecessary(s.appendName(serverConfig.Name), string(data), _etcdOps)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *etcdPositionStorageEx) Save(name string, pos mysql.Position) error {
	data, err := json.Marshal(pos)
	if err != nil {
		return err
	}

	return etcds.Save(s.appendName(name), string(data), _etcdOps)
}

func (s *etcdPositionStorageEx) Get(name string) (mysql.Position, error) {
	var entity mysql.Position

	data, _, err := etcds.Get(s.appendName(name), _etcdOps)
	if err != nil {
		return entity, err
	}

	err = json.Unmarshal(data, &entity)

	return entity, err
}
