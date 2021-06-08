package main

import (
	"flag"
	"github.com/asim/mq/go/client"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	servers = flag.String("servers", "localhost:8081", "Comma separated list of MQ servers")
)

func main() {
	flag.Parse()

	c := client.New(
		client.WithServers(strings.Split(*servers, ",")...),
	)
	tick := time.NewTicker(time.Second)
	source := `bar-`
	var count int = 0

	for _ = range tick.C {
		//count++
		text := source + strconv.Itoa(count)
		if err := c.Publish(0,"foo", []byte(text)); err != nil {
			log.Println(err)
			break
		}
		if count == 1 {
			break
		}
	}
}
