package check_server

import (
	"basic"
	"basic/tool/other"
	"errors"
	"fmt"
	"net"
	"time"
)

type CheckServer struct{}

func GetInstance() basic.Component {
	return &CheckServer{}
}

func (c *CheckServer) GetName() string {
	return "check_server"
}

func (c *CheckServer) GetDescribe() string {
	return "通过IP地址和端口号，判断应用服务是否正常或网络是否连通"
}

func (c *CheckServer) Register(globalContext *basic.Context) *basic.ComponentMeta {
	p1 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-h",
		StandardName: "host",
		Required:     true,
		CheckMethod: func(s string) error {
			if !othertool.CheckIp(s) {
				return errors.New("IP不合法")
			}
			return nil
		},
		Describe: "IP地址,支持ipv4,例:192.168.0.1",
	}
	p2 := basic.Parameter{
		ParamType:    basic.INT,
		CommandName:  "-p",
		StandardName: "port",
		Required:     true,
		CheckMethod: func(s string) error {
			if !othertool.CheckPortByString(s) {
				return errors.New("端口不合法")
			}
			return nil
		},
		Describe: "端口号,0 ~ 65535之间",
	}
	return &basic.ComponentMeta{
		ComponentType: basic.EXECUTE,
		Params:        []basic.Parameter{p1, p2},
		Component:     c,
	}
}

func (c *CheckServer) Start(globalContext *basic.Context) error {
	return nil
}

func (c *CheckServer) Do(params map[string]any) (resp []byte) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", params["host"].(string), params["port"].(int)), 3*time.Second)
	if err != nil {
		return []byte("服务不存在或网络不通")
	}
	defer conn.Close()
	return []byte("网络连通正常")
}
