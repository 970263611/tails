package web_server

import (
	"basic"
	"basic/tool/net"
	othertool "basic/tool/other"
	"errors"
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
		Required:     false,
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

func handler1(req map[string]any) (resp []byte) {
	param, ok := req["params"].(string)
	if !ok {
		return []byte("参数param不能为空")
	}
	commands := strings.Split(param, " ")
	if len(commands) <= 0 {
		return []byte("参数param不能为空")
	}
	commands = append([]string{"main.go"}, commands...)
	return basic.Servlet(commands, false)
}
