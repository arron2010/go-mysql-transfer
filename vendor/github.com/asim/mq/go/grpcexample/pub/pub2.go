package main

import (
	"fmt"
	mqgrpc "github.com/asim/mq/go/client/grpc"
	"github.com/asim/mq/logs"
)

func main() {
	logs.Initialize(nil)
	c := mqgrpc.New()
	resp, err := c.Ping("foo")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
	//text := "Hello"
	//if err := c.Publish("foo", []byte(text)); err != nil {
	//	fmt.Println(err)
	//}
}
