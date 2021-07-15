package main

import (
	"encoding/binary"
	"fmt"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"math/rand"
	"time"
)

func writeRegisterSlaveCommand(c *client.Conn, hostname string, user string, password string, serverID uint32, port uint16) error {
	c.ResetSequence()

	// This should be the name of slave host not the host we are connecting to.
	data := make([]byte, 4+1+4+1+len(hostname)+1+len(user)+1+len(password)+2+4+4)
	pos := 4

	data[pos] = mysql.COM_REGISTER_SLAVE
	pos++

	binary.LittleEndian.PutUint32(data[pos:], serverID)
	pos += 4

	// This should be the name of slave hostname not the host we are connecting to.
	data[pos] = uint8(len(hostname))
	pos++
	n := copy(data[pos:], hostname)
	pos += n

	data[pos] = uint8(len(user))
	pos++
	n = copy(data[pos:], user)
	pos += n

	data[pos] = uint8(len(password))
	pos++
	n = copy(data[pos:], password)
	pos += n

	binary.LittleEndian.PutUint16(data[pos:], port)
	pos += 2

	//replication rank, not used
	binary.LittleEndian.PutUint32(data[pos:], 0)
	pos += 4

	// master ID, 0 is OK
	binary.LittleEndian.PutUint32(data[pos:], 0)

	return c.WritePacket(data)
}

func onClientTest3() {
	//conn, _ := client.Connect("127.0.0.1:3306", "root", "12345", "eseap")
	cfg := replication.BinlogSyncerConfig{}
	cfg.ServerID = 1250
	cfg.Flavor = "mysql"
	cfg.Host = "127.0.0.1"
	cfg.Port = 3306
	cfg.User = "root"

	cfg.Password = "12345"
	cfg.Charset = "utf8"
	b := replication.NewBinlogSyncer(cfg)
	err := b.RegisterSlave()
	if err != nil {
		fmt.Println(err)
	}
	///set sql_log_bin=0;
	conn := b.GetConn()
	_, err = conn.Execute("set sql_log_bin=0;")
	if err != nil {
		fmt.Println(err)
	}
	id := uint32(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
	// Insert
	_, err = conn.Execute(fmt.Sprintf(`insert into eseap.t_user values (%d, "AAA3","BBB",10,9.85,"2020-09-24 23:33:20")`, id))
	if err != nil {
		fmt.Println(err)
	}
	conn.Close()
}
func onClientTest() {
	conn, _ := client.Connect("127.0.0.1:3306", "root", "12345", "eseap")
	//conn, _ := client.Connect("rm-2zei6e64c1k486wp18o.mysql.rds.aliyuncs.com:3306",
	//	"admin_test", "nihao123!", "uat_orderdb")
	//`uat_orderdb` eseap
	conn.Ping()
	id := uint32(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
	// Insert
	r, _ := conn.Execute(fmt.Sprintf(`insert into t_user values (%d, "AAA","BBB",10,9.85,"2020-09-24 23:33:20")`, id))

	//r, err:= conn.Execute(fmt.Sprintf(`delete from t_user where  id=%d`, id))
	//r, _ := conn.Execute(`insert into t_user values (NULL, "AA2","BB2")`)
	//r, _ := conn.Execute(`update t_user set name="AAA5" where id=1`)
	// Get last insert id
	//println(r.InsertId)
	// Or affected rows count
	//if err != nil{
	//	fmt.Println(err)
	//}

	println(r.AffectedRows)

	defer r.Close()
}

func onClientTest2() {
	conn, err := client.Connect("rm-2ze634l8642j077f8.mysql.rds.aliyuncs.com:3306", "uat_paydb", "BJtydic_456",
		"uat_paydb")
	//conn, _ := client.Connect("rm-2zei6e64c1k486wp18o.mysql.rds.aliyuncs.com:3306",
	//	"admin_test", "nihao123!", "uat_orderdb")
	//`uat_orderdb` eseap
	if err != nil {
		fmt.Println(err)
	}
	conn.Ping()
	id := uint32(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(9000))
	// Insert
	r, _ := conn.Execute(fmt.Sprintf(`insert into t_user values (%d, "CCC","DDD")`, id))

	//r, _ := conn.Execute(`insert into t_user values (NULL, "AA2","BB2")`)
	//r, _ := conn.Execute(`update t_user set name="AAA5" where id=1`)
	// Get last insert id
	println(r.InsertId)
	// Or affected rows count
	println(r.AffectedRows)

	defer r.Close()
}

func main() {
	//onClientTest2()
	onClientTest2()
}
