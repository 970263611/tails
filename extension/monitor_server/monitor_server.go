package monitor_server

import (
	"basic"
	othertool "basic/tool/other"
	"encoding/json"
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
	command.AddParameters(basic.STRING, "-h", "monitor.server.ip", "host", true, func(s string) error {
		if !othertool.CheckIp(s) {
			return errors.New("监控服务ip不合法")
		}
		return nil
	}, "监控服务的主机地址")
	command.AddParameters(basic.INT, "-p", "monitor.server.port", "port", true, func(s string) error {
		if !othertool.CheckPortByString(s) {
			return errors.New("监控服务port不合法")
		}
		return nil
	}, "监控服务的端口")
	command.AddParameters(basic.STRING, "-u", "monitor.server.username", "username", true, nil, "监控服务的登录用户名")
	command.AddParameters(basic.STRING, "-w", "monitor.server.password", "password", true, nil, "监控服务的登录密码")
	return command
}
func (r *MonitorServer) Start(globalContext *basic.Context) error {
	return nil
}

func (r *MonitorServer) Do(params map[string]any) (resp []byte) {
	res := &findResult{
		result:   &result{},
		host:     params["host"].(string),
		port:     params["port"].(int),
		username: params["username"].(string),
		password: params["password"].(string),
		c:        make(chan int),
	}
	res.urlPrefix = fmt.Sprintf("%s://%s:%d", "http", res.host, res.port)
	token, err := res.login()
	if err != nil {
		return nil
	}
	res.token = token
	go res.a1()
	go res.a2()
	go res.a3()
	go res.a4()
	go res.a5()
	go res.a6()
	go res.a10()
	after := time.After(3 * time.Second)
	for i := 8; i > 0; i-- {
		select {
		case <-res.c:
			//接口3s内执行成功
		case <-after:
			break
		}
	}
	jsonData, err := json.Marshal(res.result)
	if err != nil {
		fmt.Println("转换为JSON字符串时出错：", err)
		return
	}
	//var str string
	////str += a1 + "\r\n"
	return []byte(jsonData)
}
