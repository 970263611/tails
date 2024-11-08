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
	return "通过ip地址和端口，判断应用服务是否正常"
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
		Describe: "服务的ip地址",
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
		Describe: "",
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
		return []byte("服务不存在")
	}
	defer conn.Close()
	return []byte("服务连接成功")
}
