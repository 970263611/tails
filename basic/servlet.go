package basic

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"strings"
)

/*
*
组件分发
*/
func Servlet(commands []string) []byte {
	maps, err := commandToMap(commands)
	if err != nil {
		log.Errorf("commandToMap error : %v", err)
		return []byte(err.Error())
	}
	key, ok := maps["commands_key"]
	log.Info("call component '" + key + "'")
	if !ok {
		log.Errorf("missing component name : %v", key)
		return []byte("missing component name")
	}
	c, ok := globalContext.Components[key]
	if !ok {
		result := "component '" + key + "' not found"
		log.Error(result)
		return []byte(result)
	}
	params, err := c.Check(maps)
	if err != nil {
		log.Errorf("param check error : %v", err)
		return []byte(err.Error())
	}
	return c.Do(params)
}

/*
*
命令行转map
*/
func commandToMap(commands []string) (map[string]string, error) {
	//go run main.go component_key 至少应该有两个参数
	if len(commands) < 2 {
		return nil, errors.New("parameter error")
	}
	maps := make(map[string]string)
	//第1个参数是 office组件路径
	maps["office_path"] = commands[0]
	//第2个参数是 组件key
	commandKey := commands[1]
	maps["commandKey"] = commandKey
	//第3个到最后一个参数为 component入参，少于3个参数肯定就是没入参,只有大于等于三个参数，才存在入参
	var params []string
	if len(commands) > 2 {
		params = commands[2:len(commands)]
	}
	for i := 0; i < len(params); i++ {
		str := params[i]
		if len(str) > 1 && strings.HasPrefix(str, "-") {
			p := str[1:]
			if len(p) == 1 {
				if FindParameterType(commandKey, str) != NO_VALUE {
					i++
					maps[str] = params[i]
				}
			} else {
				for m := 0; m < len(p); m++ {
					maps["-"+p[m:m+1]] = ""
				}
			}
		} else {
			return nil, errors.New("params error")
		}
	}
	return maps, nil
}