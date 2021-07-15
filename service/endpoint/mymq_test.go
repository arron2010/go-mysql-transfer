package endpoint

import (
	"fmt"
	"strconv"
	"testing"
)

func TestMyMQEndpoint_Ping(t *testing.T) {
	a := 0x10000000
	s := strconv.FormatInt(1616686418, 16)
	//
	fmt.Println(a)
	fmt.Println(s)
}
