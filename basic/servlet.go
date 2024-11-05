package basic

import (
	othertool "basic/tool/other"
	"errors"
	log "github.com/sirupsen/logrus"
)

/*
*
组件分发
*/
func Servlet(commands []string) []byte {
	//解析命令行为map
	maps, err := commandsToMap(commands)
	if err != nil {
		log.Errorf("commandToMap error : %v", err)
		return []byte(err.Error())
	}
	//如果存在--help，则该参数优先级最高
	helpKey, ok := maps["--help"]
	if ok {
		return Help(helpKey)
	}
	key, ok := maps["componentKey"]
	if !ok {
		log.Errorf("missing component name : %v", key)
		return []byte("missing component name")
	}
	log.Info("call component '" + key + "'")
	//找到对应组件
	c, ok := globalContext.Components[key]
	if !ok {
		result := "component '" + key + "' not found"
		log.Error(result)
		return []byte(result)
	}
	//参数验证
	params, err := c.Check(maps)
	if err != nil {
		log.Errorf("param check error : %v", err)
		return []byte(err.Error())
	}
	//组件功能执行
	return c.Do(params)
}

/*
*
命令行转map
*/
func commandsToMap(commands []string) (map[string]string, error) {
	//go run main.go componentKey 至少应该有两个参数
	maps := make(map[string]string)
	//第1个参数是 office组件路径
	maps["officePath"] = commands[0]
	if len(commands) < 2 {
		maps["--help"] = "base"
		return maps, nil
	}
	//第2个参数是 组件componentkey
	componentKey := commands[1]
	maps["componentKey"] = componentKey
	//第3个到最后一个参数为 component入参，少于3个参数肯定就是没入参,只有大于等于三个参数，才存在入参
	var params []string
	if len(commands) > 2 {
		params = commands[2:len(commands)]
	}
	//如果有--help则直接返回帮助
	if commands[0] == "--help" {
		maps["--help"] = componentKey
		return maps, nil
	}
	//入参解析
	for i := 0; i < len(params); i++ {
		str := params[i]
		//参数必定以-开头，后面跟若干字母
		if othertool.RegularValidate(str, othertool.REGEX_DASE_AND_WORD) {
			p := str[1:]
			//只有单一参数才能跟值
			if len(p) == 1 {
				if FindParameterType(componentKey, str) != NO_VALUE {
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
