package basic

import (
	othertool "basic/tool/other"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

/*
*
组件分发
*/
func Servlet(commands []string, isSystem bool) []byte {
	//解析命令行为map
	maps, err := commandsToMap(commands)
	if err != nil {
		msg := fmt.Sprintf("入参解析失败 ： %v", err)
		log.Error(msg)
		return []byte(err.Error())
	}
	bytes, flag := systemParamExec(maps)
	if flag {
		return bytes
	}
	key, ok := maps["componentKey"]
	if !ok {
		msg := "入参未指定组件名称"
		log.Error(msg)
		return []byte(msg)
	}
	//找到对应组件
	c := FindComponent(key, isSystem)
	if c == nil {
		msg := fmt.Sprintf("组件 %v 不存在", key)
		log.Error(msg)
		return []byte(msg)
	}
	//参数验证
	params, err := c.Check(maps)
	if err != nil {
		msg := fmt.Sprintf("参数验证失败 : %v", err)
		log.Error(msg)
		return []byte(msg)
	}
	//组件功能执行
	return c.Do(params)
}

/*
*
系统参数列表
*/
var systemParam = []Parameter{
	Parameter{
		ParamType:    NO_VALUE,
		CommandName:  "--help",
		StandardName: "",
		Required:     false,
		Describe:     "帮助",
	},
	Parameter{
		ParamType:    STRING,
		CommandName:  "--path",
		StandardName: "",
		Required:     false,
		CheckMethod: func(s string) error {
			if s == "" {
				return errors.New("配置文件路径为空")
			}
			if !othertool.FileExists(s) {
				msg := "配置文件不存在"
				log.Error(msg)
				return errors.New(msg)
			}
			if !othertool.IsFileReadable(s) {
				msg := "配置文件不可读"
				log.Error(msg)
				return errors.New(msg)
			}
			return nil
		},
		Describe: "帮助",
	},
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
			if FindParameterType(componentKey, str) != NO_VALUE {
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
			return nil, errors.New("参数错误")
		}
	}
	return maps, nil
}

/*
*
调用Register方法
*/
func Assemble(list []Component) {
	for _, v := range list {
		cm := v.Register(globalContext)
		globalContext.Components[cm.Key] = cm
	}
}

func FindParameterType(componentKey string, commandName string) ParamType {
	for _, value := range systemParam {
		if value.CommandName == commandName {
			return value.ParamType
		}
	}
	if componentKey != "" {
		componentMeta, ok := globalContext.Components[componentKey]
		if ok {
			for _, value := range componentMeta.Params {
				if value.CommandName == commandName {
					return value.ParamType
				}
			}
		}
	}

	return -1
}

func systemParamExec(maps map[string]string) ([]byte, bool) {
	val, ok := maps["--help"]
	if ok {
		return Help(val)
	}
	val, ok = maps["--path"]
	if ok {
		return LoadConfig(val)
	}
	return nil, false
}

/*
*
方法帮助
*/
func Help(key string) ([]byte, bool) {
	linewordnum := 100
	var msg string
	if key == "base" {
		//./root 组件名 组件参数列表
		components := globalContext.Components
		msg = "命令格式:\r\n"
		msg += "  组件调用: ./root compoment_key -params\r\n"
		msg += "  组件帮助: ./root compoment_key --help\r\n"
		msg += "组件列表:"
		for _, component := range components {
			msg += "\r\n"
			msg += "  " + component.Key + "\r\n"
			chinese := othertool.SplitByChinese(component.Describe, linewordnum)
			for i := 0; i < len(chinese); i++ {
				if i == 0 {
					msg += "      " + chinese[i] + "\r\n"
				} else {
					msg += "    " + chinese[i] + "\r\n"
				}
			}
		}
	} else {
		components := FindComponent(key, false)
		if components == nil {
			msg = fmt.Sprintf("组件 %v 不存在，'./root --hlep' 查看组件列表", key)
		} else {
			params := components.Params
			msg = fmt.Sprintf("%v 参入如下:", key)
			if params != nil {
				for _, value := range components.Params {
					msg += "\r\n"
					msg += "  " + value.CommandName + "\r\n"
					msg += "    " + value.Describe
				}
			}
		}
	}
	return []byte(msg), true
}
