package main

import (
	"fmt"
	"github.com/siddontang/go-mysql/client"
	"math/rand"
	"time"
)

func onClientTest() {
	conn, _ := client.Connect("127.0.0.1:3306", "root", "12345", "eseap")
	//conn, _ := client.Connect("rm-2zei6e64c1k486wp18o.mysql.rds.aliyuncs.com:3306",
	//	"admin_test", "nihao123!", "uat_orderdb")
	//`uat_orderdb` eseap
	conn.Ping()
	id := uint32(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
	// Insert
	r, _ := conn.Execute(fmt.Sprintf(`insert into t_user values (%d, "AAA","BBB")`, id))

	//r, _ := conn.Execute(`insert into t_user values (NULL, "AA2","BB2")`)
	//r, _ := conn.Execute(`update t_user set name="AAA5" where id=1`)
	// Get last insert id
	println(r.InsertId)
	// Or affected rows count
	println(r.AffectedRows)

	defer r.Close()
}

func onClientTest2() {
	conn, _ := client.Connect("rm-2ze634l8642j077f8.mysql.rds.aliyuncs.com:3306", "uat_paydb", "BJtydic_456", "uat_paydb")
	//conn, _ := client.Connect("rm-2zei6e64c1k486wp18o.mysql.rds.aliyuncs.com:3306",
	//	"admin_test", "nihao123!", "uat_orderdb")
	//`uat_orderdb` eseap
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
	onClientTest2()
	onClientTest()
}
