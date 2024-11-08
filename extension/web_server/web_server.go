package web_server

import "basic"

type WebServer struct{}

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

	return nil
}

func (r *WebServer) Start(globalContext *basic.Context) error {
	return nil
}

func (r *WebServer) Do(params map[string]any) (resp []byte) {
	return nil
}
