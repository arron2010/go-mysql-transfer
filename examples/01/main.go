package main

import (
	"fmt"
	"github.com/siddontang/go-mysql/canal"
	"os"
	"os/signal"
	"syscall"
)

type MyEventHandler struct {
	canal.DummyEventHandler
}

func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
	fmt.Printf(">>>>     %s %v\n", e.Action, e.Rows)
	return nil
}

func (h *MyEventHandler) String() string {
	return "MyEventHandler"
}
func wait() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Kill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sin := <-s
	fmt.Printf("application stoped，signal: %s \n", sin.String())
}

func onServerTest() {
	cfg := canal.NewDefaultConfig()
	//cfg.Addr = "rm-2zei6e64c1k486wp18o.mysql.rds.aliyuncs.com:3306"
	//cfg.User = "admin_test"
	//cfg.Password="nihao123!"

	cfg.Addr = "127.0.0.1:3306"
	cfg.User = "root"
	cfg.Password = "12345"

	// We only care table canal_test in test db
	cfg.Dump.TableDB = "eseap"
	//cfg.Dump.Databases=[]string{"eseap"}
	cfg.Dump.Tables = []string{"t_user"}

	c, err := canal.NewCanal(cfg)
	if err != nil {
		fmt.Println("onServerTest-->", err)
	}

	// Register a handler to handle RowsEvent
	c.SetEventHandler(&MyEventHandler{})

	// Start canal
	c.Run()
}
func main() {
	onServerTest()
	wait()
}
