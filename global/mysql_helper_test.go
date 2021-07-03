package global

import (
	"fmt"

	"testing"
)

func TestGetLastBinlog(t *testing.T) {
	config := &ServerConfig{}
	config.Addr = "127.0.0.1:3306"
	config.User = "root"
	config.Password = "12345"
	helper := NewMySQLHelper(config)
	file, pos, err := helper.GetLastBinlog()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(file, " ", pos)
}
