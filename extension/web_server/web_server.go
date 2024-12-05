package web_server

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/utils"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type WebServer struct {
	iface.Context
}

var webServer *WebServer

func GetInstance(globalContext iface.Context) iface.Component {
	webServer = &WebServer{globalContext}
	return webServer
}

func (w *WebServer) GetName() string {
	return "web_server"
}

func (r *WebServer) GetDescribe() string {
	return "启动web服务\n例：web_server -p 17001"
}

func (r *WebServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.INT, cons.LOWER_P, "port", "port", true,
		func(s string) error {
			if !utils.CheckPortByString(s) {
				return errors.New("端口不合法")
			}
			return nil
		}, "发布web服务的端口号")
}

func (r *WebServer) Do(params map[string]any) (resp []byte) {
	port := params["port"].(int)
	r.Context.GetConfig().SetDefault("web_server.port", port)
	portStr := strconv.Itoa(port)
	log.Info("Web will start at " + portStr)
	err := http.ListenAndServe(":"+portStr, nil)
	return []byte("服务启动失败:" + err.Error())
}
