package tcp_num

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	params := make(map[string]any)
	//params["tcp_port"] = 6379
	//params["tcp_established"] = ""
	params["tcp_ip"] = "127.0.0.1"
	ins := GetInstance()
	bytes := ins.Do(params)
	fmt.Println(bytes)
	fmt.Println(string(bytes))
}
