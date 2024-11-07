package main

import (
	"basic"
	"basic/log_config"
	"basic/onload"
	_ "basic/onload"
	othertool "basic/tool/other"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	//初始化日志配置
	config := log_config.NewLogConfig()
	config.OutType = 0
	log_config.Init(config)
}

func main() {
	args := os.Args
	//args := []string{"main.go", "check_server", "--help"}
	//解析命令行为map
	maps, err := commandsToMap(args)
	if err != nil {
		msg := fmt.Sprintf("入参解析失败 ： %v", err)
		log.Error(msg)
		fmt.Println(string(msg))
	}
	onload.Init()
	//调用component，并打印
	bytes := basic.Servlet(maps, false)
	fmt.Println(string(bytes))
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
		params = commands[2:]
	}
	//入参解析
	for i := 0; i < len(params); i++ {
		str := params[i]
		//系统参数必定以--开头，后面跟若干字母
		if othertool.RegularValidate(str, othertool.REGEX_2DASE_AND_WORD) {
			if basic.FindParameterType(componentKey, str) != basic.NO_VALUE {
				i++
				maps[str] = params[i]
			} else {
				maps[str] = componentKey
			}
			//组件参数必定以-开头，后面跟若干字母
		} else if othertool.RegularValidate(str, othertool.REGEX_DASE_AND_WORD) {
			p := str[1:]
			//只有单一参数才能跟值
			if len(p) == 1 {
				if basic.FindParameterType(componentKey, str) != basic.NO_VALUE {
					i++
					maps[str] = params[i]
				}
			} else {
				for m := 0; m < len(p); m++ {
					maps["-"+p[m:m+1]] = ""
				}
			}
		} else {
			return nil, errors.New("参数错误")
		}
	}
	return maps, nil
}
