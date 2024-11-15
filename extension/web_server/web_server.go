package web_server

import (
	"basic"
	"basic/tool/net"
	othertool "basic/tool/other"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type WebServer struct {
	*basic.Context
}

func GetInstance() *WebServer {
	return &WebServer{}
}

func (w *WebServer) GetName() string {
	return "web_server"
}

func (r *WebServer) GetDescribe() string {
	return "启动web服务"
}

func (r *WebServer) Register(globalContext *basic.Context) *basic.ComponentMeta {
	r.Context = globalContext
	p1 := basic.Parameter{
		ParamType:    basic.INT,
		CommandName:  "-p",
		StandardName: "port",
		ConfigName:   "web_server.port",
		Required:     true,
		CheckMethod: func(s string) error {
			if !othertool.CheckPortByString(s) {
				return errors.New("端口不合法")
			}
			return nil
		},
		Describe: "发布web服务的端口号",
	}
	return &basic.ComponentMeta{
		ComponentType: basic.EXECUTE,
		Params:        []basic.Parameter{p1},
		Component:     r,
	}
}

func (r *WebServer) Start(globalContext *basic.Context) error {
	return nil
}

func (r *WebServer) Do(params map[string]any) (resp []byte) {
	port, ok := params["port"].(int)
	if !ok {
		port = r.Config.GetInt("web_server.server.port")
	}
	if port == 0 {
		return []byte("未设置服务端口号")
	}
	handlers := make(map[string]func(req map[string]any) (resp []byte))
	handlers["/do"] = handler1
	err := net.Web(port, handlers)
	return []byte("服务启动失败:" + err.Error())
}

/*
*
web请求处理逻辑
*/
func handler1(req map[string]any) []byte {
	param, ok := req["params"].(string)
	if !ok {
		return []byte("参数param不能为空")
	}
	commands := othertool.SplitString(param)
	if len(commands) <= 0 {
		return []byte("参数param不能为空")
	}
	resp, b := forward(commands)
	if b {
		return resp
	}
	commands = append([]string{"main.go"}, commands...)
	return basic.Servlet(commands, false)
}

/*
*
web请求转发
*/
func forward(args []string) ([]byte, bool) {
	var addr, params string
	var flag bool
	//判断是否需要转发，并拼接转发参数
	for i := 0; i < len(args); i++ {
		if args[i] == "--addr" {
			flag = true
			i++
			if i < len(args) {
				addr = args[i]
			}
			if addr == "" {
				msg := fmt.Sprintf("请求转发的ip端口为空")
				log.Error(msg)
				return []byte(msg), true
			}
		} else {
			params = args[i] + " "
		}
	}
	if !flag {
		return nil, flag
	}
	//校验转发地址是否合法
	msg := fmt.Sprintf("地址不合法")
	arr := strings.Split(addr, ":")
	if len(arr) != 2 {
		log.Error(msg)
		return []byte(msg), flag
	}
	ipflag := othertool.CheckIp(arr[0])
	portflag := othertool.CheckPortByString(arr[1])
	if !ipflag || !portflag {
		log.Error(msg)
		return []byte(msg), flag
	}
	//转发请求并返回结果
	uri := fmt.Sprintf("http://%s/do", addr)
	resp, err := net.PostRespString(uri, map[string]string{
		"params": params,
	}, nil)
	if err != nil {
		return []byte(err.Error()), flag
	} else {
		return []byte(resp), flag
	}
}
