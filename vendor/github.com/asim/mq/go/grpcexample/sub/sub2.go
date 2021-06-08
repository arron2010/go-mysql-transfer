package main

import (
	"fmt"
	mqgrpc "github.com/asim/mq/go/client/grpc"
)

func main() {
	c := mqgrpc.New()
	ch, err := c.Subscribe("foo")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		select {
		case e := <-ch:
			fmt.Println(string(e))
		}
	}
}
