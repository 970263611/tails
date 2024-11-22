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
	return "启动web服务,例:web_server -p 17001"
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
	handlers := make(map[string]func(req map[string]any) (resp []byte))
	r.Context.GetConfig().SetDefault("web_svc.port", port)
	handlers["/do"] = r.handler1
	err := net.Web(port, handlers)
	return []byte("服务启动失败:" + err.Error())
}

/*
*
web请求处理逻辑
*/
func (r *WebServer) handler1(req map[string]any) []byte {
	defer r.DelCache()
	param, ok := req["params"].(string)
	_, isSystem := req["isSystem"]
	if !ok {
		return msgFormat("", "参数param不能为空", isSystem)
	}
	commands := utils.SplitString(param)
	if len(commands) <= 0 {
		return msgFormat("", "参数param不能为空", isSystem)
	}
	//解析入参中的系统参数
	args, err := r.LoadSystemParams(commands)
	if err != nil {
		return msgFormat("", err.Error(), isSystem)
	}
	data, err := r.Servlet(args, false)
	if err != nil {
		return msgFormat("", err.Error(), isSystem)
	} else {
		return msgFormat(string(data), "", isSystem)
	}
}

func msgFormat(data string, errmsg string, isSystem bool) []byte {
	if isSystem {
		rmap := map[string]string{}
		if errmsg != "" {
			rmap[cons.RESULT_CODE] = cons.RESULT_ERROR
			rmap[cons.RESULT_DATA] = errmsg
		} else {
			rmap[cons.RESULT_CODE] = cons.RESULT_SUCCESS
			rmap[cons.RESULT_DATA] = data
		}
		jsonData, err := json.Marshal(rmap)
		if err != nil {
			log.Errorf("转换报文的map:%v", rmap)
			log.Errorf("格式化报文错误:%v", err)
			return []byte(err.Error())
		}
		return jsonData
	} else {
		if errmsg != "" {
			return []byte(errmsg)
		} else {
			return []byte(data)
		}
	}
}
