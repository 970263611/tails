package basic

import (
	"basic/tool/net"
	"basic/tool/utils"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

const (
	COMPONENT_KEY string = "componentKey"
	WEB_KEY       string = "web_server"
	NEEDHELP      string = "--help"
	GID           string = "globalID"
)

/*
*
组件分发
*/
func Servlet(commands []string, isSystem bool) []byte {
	defer DelCache()
	commands = setGID(commands)
	//参数中携带 --addr ip:port或域名 时进行请求转发
	resp, b := forward(commands)
	if b {
		return resp
	}
	return execute(commands, isSystem)
}

/*
*
命令执行
*/
func execute(commands []string, isSystem bool) []byte {
	//解析命令行为map
	maps, err := commandsToMap(commands)
	if err != nil {
		msg := fmt.Sprintf("入参解析失败 ： %v", err)
		log.Error(msg)
		return []byte(msg)
	}
	//获取组件帮助信息
	if key, ok := maps[NEEDHELP]; ok {
		return help(key)
	}
	//配置文件配置注入参数中
	addConfigToMap(maps)
	key, ok := maps[COMPONENT_KEY]
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
	if key == WEB_KEY {
		errStart := start()
		if errStart != nil {
			return []byte(errStart.Error())
		}
	}
	//组件功能执行
	log.Infof("组件执行开始:[%v],组件入参:%v", key, commands)
	res := c.Do(params)
	if res != nil {
		log.Infof("组件执行完成:[%v],组件执行结果大小:[%d byte]", key, len(res))
	}
	log.Debugf("组件执行结果:%v", string(res))
	return res
}

/*
*
web请求转发
*/
func forward(commands []string) ([]byte, bool) {
	var addr, params string
	var flag bool
	//判断是否需要转发，并拼接转发参数
	for i := 0; i < len(commands); i++ {
		if commands[i] == "--addr" {
			flag = true
			i++
			if i < len(commands) {
				addr = commands[i]
			}
			if addr == "" {
				msg := fmt.Sprintf("请求转发的ip端口为空")
				log.Error(msg)
				return []byte(msg), true
			}
		} else {
			params += commands[i] + " "
		}
	}
	if !flag {
		return nil, flag
	}
	//校验转发地址是否合法
	if !utils.CheckAddr(addr) {
		msg := fmt.Sprintf("地址不合法")
		log.Error(msg)
		return []byte(msg), flag
	}
	log.Infof("请求转发，转发地址:[%s],转发参数:[%s]", addr, params)
	//插入全局业务跟踪号
	gid, ok := GetCache(GID)
	if ok {
		id := fmt.Sprintf("%v", gid)
		params += " --gid "
		params += id
	}
	//转发请求并返回结果
	uri := fmt.Sprintf("http://%s/do", addr)
	resp, err := net.PostRespString(uri, map[string]string{
		"params": params,
	}, nil)
	if err != nil {
		log.Errorf("转发请求 %v 错误，错误原因 : &v", uri, err)
		return []byte(err.Error()), flag
	} else {
		return []byte(resp), flag
	}
}

/*
*
设置全局ID
*/
func setGID(commands []string) []string {
	_, ok := GetCache(GID)
	if ok {
		return commands
	}
	var gid string
	var commands2 []string
	//获取全局gid参数
	for i := 0; i < len(commands); i++ {
		if commands[i] == "--gid" {
			i++
			if i < len(commands) {
				gid = commands[i]
			}
		} else {
			commands2 = append(commands2, commands[i])
		}
	}
	if gid == "" {
		gid, _ = utils.GenerateUUID()
	}
	SetCache(GID, gid)
	return commands2
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
		Describe:     "查看系统或组件帮助",
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
			if !utils.FileExists(s) {
				msg := "配置文件不存在"
				log.Error(msg)
				return errors.New(msg)
			}
			if !utils.IsFileReadable(s) {
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
方法帮助
*/
func help(key string) []byte {
	return []byte(globalContext.FindHelp(key))
}

/*
*
start执行
*/
func start() error {
	err := globalContext.Start()
	if err != nil {
		return err
	}
	return nil
}
