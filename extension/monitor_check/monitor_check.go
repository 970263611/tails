package monitor_check

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

type MonitorServer struct {
	iface.Context
}

func GetInstance(globalContext iface.Context) iface.Component {
	return &MonitorServer{globalContext}
}

func (r *MonitorServer) GetName() string {
	return "monitor_check"
}

func (r *MonitorServer) GetDescribe() string {
	return "应急一键检查指标\n例：monitor_check -h 127.0.0.1 -p 8080 -u sysadmin -w Psbc@2023"
}

func (r *MonitorServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, cons.LOWER_H, "ip", "host", true, func(s string) error {
		if !utils.CheckIp(s) {
			return errors.New("监控服务ip不合法")
		}
		return nil
	}, "监控服务的主机地址")
	cm.AddParameters(cons.INT, cons.LOWER_P, "port", "port", true, func(s string) error {
		if !utils.CheckPortByString(s) {
			return errors.New("监控服务port不合法")
		}
		return nil
	}, "监控服务的端口")
	cm.AddParameters(cons.STRING, cons.LOWER_U, "username", "username", true, nil, "监控服务的登录用户名")
	cm.AddParameters(cons.STRING, cons.LOWER_W, "password", "password", true, nil, "监控服务的登录密码")
	cm.AddParameters(cons.STRING, cons.UPPER_V, "version", "version", false, nil, "监控服务的登录版本1无需验证码(默认) 2需要验证码")
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
	version, flag := params["version"].(string)
	if !flag {
		version = "1"
	}
	if version == "2" {
		//验证码登录
		res.version = "2"
		// 获取登录页面验证码
		identifyCode, err := res.getIdentifyCode()
		if err != nil {
			return []byte("获取验证码失败: " + err.Error())
		}
		if identifyCode.Code == 0 {
			res.identifyCode = identifyCode.Data.Code
			res.identifyCodeUuid = identifyCode.Data.CodeUuid
		} else {
			return []byte("获取验证码失败: " + err.Error())
		}
	} else {
		//默认无需验证码登录
		res.version = "1"
	}
	//登录
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
		log.Error("转换为JSON字符串时出错", err.Error())
		return []byte("转换为JSON字符串时出错" + err.Error())
	}
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, jsonData, "", "    ")
	if error != nil {
		log.Error("JSON美化出错", err.Error())
		return []byte("JSON美化出错" + error.Error())
	}
	prettyStr := string(prettyJSON.Bytes())
	return []byte(prettyStr)
}
