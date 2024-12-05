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
	distribute(resp, result)
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
func distribute(resp http.ResponseWriter, req map[string]any) {
	defer webServer.DelCache()
	params, ok := req["params"].(string)
	if !ok {
		params = cons.SYSPARAM_HELP
	}
	commands := utils.SplitString(params)
	if len(commands) <= 0 {
		http.Error(resp, "参数params不能为空", http.StatusBadRequest)
		return
	}
	//解析入参中的系统参数
	args, err := webServer.LoadSystemParams(commands)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	webServer.SetCache(cons.WEB_RESP, resp)
	data, err := webServer.Servlet(args, false)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusAccepted)
		return
	} else {
		if data != nil {
			resp.Write(data)
		}
		return
	}
}
