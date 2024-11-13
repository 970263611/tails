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
func Servlet(args []string, isSystem bool) []byte {
	//解析命令行为map
	maps, err := commandsToMap(args)
	if err != nil {
		msg := fmt.Sprintf("入参解析失败 ： %v", err)
		log.Error(msg)
		return []byte(msg)
	}
	err = addConfigToMap(maps)
	if err != nil {
		msg := fmt.Sprintf("配置文件参数解析失败 ： %v", err)
		log.Error(msg)
		return []byte(msg)
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
	c := globalContext.FindComponent(key, isSystem)
	if c == nil {
		msg := fmt.Sprintf("组件 %v 不存在", key)
		log.Error(msg)
		return []byte(msg)
	}
	//参数验证
	params, err := c.check(maps)
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
		ParamType:    NO_VALUE,
		CommandName:  "--start",
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

func systemParamExec(maps map[string]string) ([]byte, bool) {
	val, ok := maps["--help"]
	if ok {
		return Help(val)
	}
	val, ok = maps["--start"]
	if ok {
		return Start()
	}
	return nil, false
}

/*
*
方法帮助
*/
func Help(key string) ([]byte, bool) {
	help := globalContext.FindHelp(key)
	return []byte(help), true
}

/*
*
start执行
*/
func Start() ([]byte, bool) {
	err := globalContext.Start()
	if err != nil {
		return []byte(err.Error()), false
	}
	return nil, false
}
