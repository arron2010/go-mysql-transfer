package main

import (
	"fmt"
	"github.com/siddontang/go-mysql/client"
	"strconv"
	"time"
)

func insertData(id uint32) {
	conn, _ := client.Connect("127.0.0.1:3306", "root", "12345", "eseap")
	sql := "insert into t_user(`id`,`name`,`desc`,`a`,`b`,`c`) values(?,?,?,?,?,?)"

	tt := time.Now()
	const format = "2006-01-02 15:04:05"

	//id := uint32(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
	name := "AAA" + strconv.Itoa(int(id))
	desc := "AAA" + strconv.Itoa(int(id))
	a := id
	b := float32(id) / 100
	c := tt.Format(format)
	_, err := conn.Execute(sql, id, name, desc, a, b, c)
	if err != nil {
		fmt.Println("Error-->", err)
	}
	defer conn.Close()
}
func insertData2(id uint32) {
	conn, _ := client.Connect("127.0.0.1:3306", "root", "12345", "eseap2")
	sql := "insert into t_user3(`id`,`name`,`desc`,`a`,`b`,`c`) values(?,?,?,?,?,?)"

	tt := time.Now()
	const format = "2006-01-02 15:04:05"

	//id := uint32(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
	name := "AAA" + strconv.Itoa(int(id))
	desc := "AAA" + strconv.Itoa(int(id))
	a := id
	b := float32(id) / 100
	c := tt.Format(format)
	_, err := conn.Execute(sql, id, name, desc, a, b, c)
	if err != nil {
		fmt.Println("Error-->", err)
	}
	defer conn.Close()
}
func updateData() {
	conn1, _ := client.Connect("127.0.0.1:3306", "root", "12345", "eseap")
	defer conn1.Close()
	sql1 := `update t_user t set t.desc="KKK" where t.id=12`
	_, err := conn1.Execute(sql1)
	if err != nil {
		fmt.Println("eseap Error-->", err)
	}

	conn2, _ := client.Connect("127.0.0.1:3306", "root", "12345", "eseap2")
	defer conn2.Close()
	sql2 := `update t_user3 t set t.desc="MMM" where t.id=12`
	_, err = conn2.Execute(sql2)
	if err != nil {
		fmt.Println("eseap2 Error-->", err)
	}
}
func main() {

	//for i :=1; i <= 10;i++{
	//	insertData(uint32(i))
	//}
	//fmt.Println("写入数据成功！")

	//for i :=12; i <= 12;i++{
	//	insertData2(uint32(i))
	//}

	updateData()
	fmt.Println("写入数据成功！")
	//utils.WaitFor()
}
