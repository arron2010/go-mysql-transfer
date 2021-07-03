package global

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pingcap/errors"
	"strings"
)

type ConfigDAO struct {
	dsn string
}

func NewConfigDAO(addr string, user string, password string, charset string) *ConfigDAO {
	//root:12345@tcp(127.0.0.1:3306)/configdb
	dao := &ConfigDAO{}
	dao.dsn = fmt.Sprintf("%s:%s@tcp(%s)/configdb", user, password, addr)
	return dao
}

func (c *ConfigDAO) createRule(path string) *Rule {
	pathInfo := strings.Split(path, "/")
	rule := &Rule{}
	rule.Schema = pathInfo[2]
	rule.Table = pathInfo[3]
	return rule
}
func (c *ConfigDAO) CreateServerConfig() ([]*ServerConfig, error) {

	db, err := sql.Open("mysql", c.dsn)
	defer db.Close()
	if err != nil {
		return nil, errors.Trace(err)
	}

	m := make(map[string]*ServerConfig)

	rows, err := db.Query(`select t2.id, t2.name_,t2.user_,t2.password_,t2.addr_,t1.path_ FROM t_config_strategy t1,t_config_dbinfo t2 where t1.server_=t2.name_ and t1.deleted_="0"`)

	servers := make([]*ServerConfig, 0, 8)
	for rows.Next() {
		s := &ServerConfig{}
		s.RuleConfigs = make([]*Rule, 0, 0)
		err := rows.Scan(&s.SlaveID, &s.Name, &s.User, &s.Password, &s.Addr, &s.Path)
		s.SlaveID = s.SlaveID + 1000
		if err != nil {
			continue
		}
		old, ok := m[s.Name]
		if ok {
			old.RuleConfigs = append(old.RuleConfigs, c.createRule(s.Path))
		} else {
			s.RuleConfigs = append(s.RuleConfigs, c.createRule(s.Path))
			servers = append(servers, s)
			m[s.Name] = s
		}
	}
	rows.Close()
	return servers, err
}
