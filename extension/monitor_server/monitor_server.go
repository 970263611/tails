package monitor_server

import (
	"basic"
	othertool "basic/tool/other"
	"errors"
	"fmt"
	"time"
)

type MonitorServer struct {
	*basic.Context
}

func GetInstance() *MonitorServer {
	return &MonitorServer{}
}

func (w *MonitorServer) GetName() string {
	return "monitor_server"
}

func (r *MonitorServer) GetDescribe() string {
	return "monitor监控指标服务"
}

func (r *MonitorServer) Register(globalContext *basic.Context) *basic.ComponentMeta {
	command := &basic.ComponentMeta{
		Component: r,
	}
	command.AddParameters(basic.STRING, "-h", "host", true, func(s string) error {
		if !othertool.CheckIp(s) {
			return errors.New("监控服务ip不合法")
		}
		return nil
	}, "监控服务的主机地址")
	command.AddParameters(basic.INT, "-p", "port", true, func(s string) error {
		if !othertool.CheckPortByString(s) {
			return errors.New("监控服务port不合法")
		}
		return nil
	}, "监控服务的端口")
	command.AddParameters(basic.STRING, "-u", "username", true, nil, "监控服务的登录用户名")
	command.AddParameters(basic.STRING, "-w", "password", true, nil, "监控服务的登录密码")
	return command
}
func (r *MonitorServer) Start(globalContext *basic.Context) error {
	return nil
}

var urlPrefix string

func (r *MonitorServer) Do(params map[string]any) (resp []byte) {
	res := &findResult{
		result:   &result{},
		host:     params["host"].(string),
		port:     params["port"].(int),
		username: params["username"].(string),
		password: params["password"].(string),
		c:        make(chan int),
	}
	token, err := res.login()
	if err != nil {
		return nil
	}
	i := 4
	res.token = token
	go res.a1()
	go res.a2()
	go res.a3()
	go res.a4()
	after := time.After(3 * time.Second)
	for {
		select {
		case <-res.c:
			i--
			if i == 0 {
				break
			}
		case <-after:
			break
		}
	}
	str := fmt.Sprintf("%v", res.result)
	return []byte(str)
}
