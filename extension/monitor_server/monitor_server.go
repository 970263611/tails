package monitor_server

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type MonitorServer struct {
	iface.Context
}

func GetInstance(globalContext iface.Context) iface.Component {
	return &MonitorServer{globalContext}
}

func (r *MonitorServer) GetName() string {
	return "monitor_server"
}

func (r *MonitorServer) GetDescribe() string {
	return "monitor监控指标服务"
}

func (r *MonitorServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, "-h", "monitor.server.ip", "host", true, func(s string) error {
		if !utils.CheckIp(s) {
			return errors.New("监控服务ip不合法")
		}
		return nil
	}, "监控服务的主机地址")
	cm.AddParameters(cons.INT, "-p", "monitor.server.port", "port", true, func(s string) error {
		if !utils.CheckPortByString(s) {
			return errors.New("监控服务port不合法")
		}
		return nil
	}, "监控服务的端口")
	cm.AddParameters(cons.STRING, "-u", "monitor.server.username", "username", true, nil, "监控服务的登录用户名")
	cm.AddParameters(cons.STRING, "-w", "monitor.server.password", "password", true, nil, "监控服务的登录密码")
}
func (r *MonitorServer) Start() error {
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
		return []byte("登录失败: " + err.Error())
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
		return []byte("转换为JSON字符串时出错")
	}
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, jsonData, "", "    ")
	if error != nil {
		return []byte("JSON美化出错")
	}
	prettyStr := string(prettyJSON.Bytes())
	return []byte(prettyStr)
}
