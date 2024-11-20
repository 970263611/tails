package web_server

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/net"
	"basic/tool/utils"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
)

type WebServer struct {
	iface.Context
}

func GetInstance(globalContext iface.Context) iface.Component {
	return &WebServer{globalContext}
}

func (w *WebServer) GetName() string {
	return "web_server"
}

func (r *WebServer) GetDescribe() string {
	return "启动web服务"
}

func (r *WebServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.INT, cons.LOWER_P, "web_server.port", "port", true,
		func(s string) error {
			if !utils.CheckPortByString(s) {
				return errors.New("端口不合法")
			}
			return nil
		}, "发布web服务的端口号")
}

func (r *WebServer) Do(params map[string]any) (resp []byte) {
	port := params["port"].(int)
	handlers := make(map[string]func(req map[string]any) (resp []byte))
	handlers["/do"] = r.handler1
	err := net.Web(port, handlers)
	return []byte("服务启动失败:" + err.Error())
}

/*
*
web请求处理逻辑
*/
func (r *WebServer) handler1(req map[string]any) []byte {
	param, ok := req["params"].(string)
	if !ok {
		return []byte("参数param不能为空")
	}
	commands := utils.SplitString(param)
	if len(commands) <= 0 {
		return []byte("参数param不能为空")
	}
	data, err := r.Servlet(commands, false)
	_, ok = req["isSystem"]
	if ok {
		rmap := map[string]string{}
		if err != nil {
			rmap[cons.RESULT_CODE] = cons.RESULT_ERROR
			rmap[cons.RESULT_DATA] = err.Error()
		} else {
			rmap[cons.RESULT_CODE] = cons.RESULT_SUCCESS
			rmap[cons.RESULT_DATA] = string(data)
		}
		jsonData, err := json.Marshal(rmap)
		if err != nil {
			log.Errorf("JSON encoding failed:%v", err)
			return []byte(err.Error())
		}
		return jsonData
	} else {
		if err != nil {
			return []byte(err.Error())
		} else {
			return []byte(data)
		}
	}
}
