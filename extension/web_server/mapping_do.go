package web_server

import (
	cons "basic/constants"
	"basic/tool/utils"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func init() {
	http.HandleFunc("/do", do)
}

func do(resp http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	req.ParseForm()
	maps := mergeMaps(mapChange(req.Form, false), mapChange(params, true))
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Error("Error reading request body", err)
	}
	var p3 map[string]any
	if len(body) > 0 {
		err = json.Unmarshal(body, &p3)
		if err != nil {
			log.Error(resp, "Error parsing JSON", err)
		}
	}
	result := mergeMaps(maps, p3)
	msgbyte := distribute(result)
	if msgbyte != nil {
		resp.Write(msgbyte)
	}
}

func mapChange(m1 map[string][]string, query bool) map[string]any {
	changed := make(map[string]any)
	for k, v := range m1 {
		if query {
			if len(v) == 1 {
				changed[k] = v[0]
				continue
			}
		}
		changed[k] = v
	}
	return changed
}

func mergeMaps(m1 map[string]any, m2 map[string]any) map[string]any {
	merged := make(map[string]any)
	for k, v := range m1 {
		merged[k] = v
	}
	for k, v := range m2 {
		merged[k] = v
	}
	return merged
}

/*
*
web请求处理逻辑
*/
func distribute(req map[string]any) []byte {
	defer webServer.DelCache()
	params, ok := req["params"].(string)
	_, isSystem := req["isSystem"]
	if !ok {
		params = cons.SYSPARAM_HELP
	}
	commands := utils.SplitString(params)
	if len(commands) <= 0 {
		return msgFormat("", "参数params不能为空", isSystem)
	}
	//解析入参中的系统参数
	args, err := webServer.LoadSystemParams(commands)
	if err != nil {
		return msgFormat("", err.Error(), isSystem)
	}
	data, err := webServer.Servlet(args, false)
	if err != nil {
		return msgFormat("", err.Error(), isSystem)
	} else {
		return msgFormat(string(data), "", isSystem)
	}
}

/*
*
返回报文格式化
*/
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
