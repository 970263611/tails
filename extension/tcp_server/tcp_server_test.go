package tcp_server

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	params := make(map[string]any)
	//params["tcp_port"] = 6379
	//params["tcp_established"] = ""
	ins := GetInstance()
	bytes := ins.Do(params)
	fmt.Println(bytes)
	fmt.Println(string(bytes))
}