package check_server

import (
	"basic"
	"basic/tool/other"
	"errors"
	"fmt"
	"math"
	"net"
	"time"
)

type CheckServer struct{}

func GetInstance() basic.Component {
	return &CheckServer{}
}

func (c *CheckServer) GetOrder() int {
	return math.MaxInt64
}

func (c *CheckServer) Register(globalContext *basic.Context) *basic.ComponentMeta {
	p1 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-h",
		StandardName: "host",
		Required:     true,
		CheckMethod: func(s string) error {
			if !othertool.CheckIp(s) {
				errors.New("port is not valid")
			}
			return nil
		},
		Describe: "服务的ip地址",
	}
	p2 := basic.Parameter{
		ParamType:    basic.INT,
		CommandName:  "-p",
		StandardName: "host",
		Required:     true,
		CheckMethod: func(s string) error {
			if !othertool.CheckPortByString(s) {
				errors.New("port is not valid")
			}
			return nil
		},
		Describe: "",
	}
	return &basic.ComponentMeta{
		ComponentType: basic.EXECUTE,
		Key:           "check_server",
		Describe:      "服务的端口",
		Params:        []basic.Parameter{p1, p2},
		Component:     c,
	}
}

func (c *CheckServer) Do(params map[string]any) (resp []byte) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", params["host"].(string), params["port"].(int)), 3*time.Second)
	if err != nil {
		return []byte("disconnected")
	}
	defer conn.Close()
	return []byte("connected")
}
