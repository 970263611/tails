package basic

import (
	cons "basic/constants"
	"basic/tool/net"
	"basic/tool/utils"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

/*
*
组件分发
*/
func (c *Context) Servlet(commands []string, isSystem bool) ([]byte, error) {
	defer globalContext.DelCache()
	commands = setGID(commands)
	//参数中携带 --addr ip:port或域名 时进行请求转发
	resp, err, b := forward(commands)
	if b {
		return resp, err
	}
	return execute(commands, isSystem)
}

/*
*
命令执行
*/
func execute(commands []string, isSystem bool) ([]byte, error) {
	//解析命令行为map
	maps, err := commandsToMap(commands)
	if err != nil {
		msg := fmt.Sprintf("入参解析失败 ： %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}
	//获取组件帮助信息
	if key, ok := maps[cons.NEEDHELP]; ok {
		return help(key), nil
	}
	//配置文件配置注入参数中
	addConfigToMap(maps)
	key, ok := maps[cons.COMPONENT_KEY]
	if !ok {
		msg := "入参未指定组件名称"
		log.Error(msg)
		return nil, errors.New(msg)
	}
	//找到对应组件
	c := globalContext.findComponent(key, isSystem)
	if c == nil {
		msg := fmt.Sprintf("组件 %v 不存在", key)
		log.Error(msg)
		return nil, errors.New(msg)
	}
	//参数验证
	params, err := c.check(maps)
	if err != nil {
		msg := fmt.Sprintf("参数验证失败 : %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}
	//组件功能执行
	log.Infof("组件执行开始:[%v],组件入参:%v", key, commands)
	res := c.Do(params)
	if res != nil {
		log.Infof("组件执行完成:[%v],组件执行结果大小:[%d byte]", key, len(res))
	}
	log.Debugf("组件执行结果:%v", string(res))
	return res, nil
}

/*
*
web请求转发
*/
func forward(commands []string) ([]byte, error, bool) {
	var addr, params string
	var flag bool
	//判断是否需要转发，并拼接转发参数
	for i := 0; i < len(commands); i++ {
		if commands[i] == "--forward" || commands[i] == "--f" {
			flag = true
			i++
			if i < len(commands) {
				addr = commands[i]
			}
			if addr == "" {
				msg := fmt.Sprintf("请求转发的ip端口为空")
				log.Error(msg)
				return nil, errors.New(msg), true
			}
		} else {
			params += commands[i] + " "
		}
	}
	if !flag {
		return nil, nil, flag
	}
	//校验转发地址是否合法
	if !utils.CheckAddr(addr) {
		msg := fmt.Sprintf("地址不合法")
		log.Error(msg)
		return nil, errors.New(msg), flag
	}
	log.Infof("请求转发，转发地址:[%s],转发参数:[%s]", addr, params)
	//插入全局业务跟踪号
	gid, ok := globalContext.GetCache(cons.GID)
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
		return nil, err, flag
	} else {
		return []byte(resp), nil, flag
	}
}

/*
*
设置全局ID
*/
func setGID(commands []string) []string {
	_, ok := globalContext.GetCache(cons.GID)
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
	globalContext.SetCache(cons.GID, gid)
	return commands2
}

/*
*
方法帮助
*/
func help(key string) []byte {
	return []byte(globalContext.findHelp(key))
}

/*
*
命令行数组转换为maps
*/
func commandsToMap(commands []string) (map[string]string, error) {
	//go run main.go componentKey 至少应该有两个参数
	maps := make(map[string]string)
	if len(commands) == 0 {
		maps["--help"] = "base"
		return maps, nil
	}
	//第1个参数是 组件componentkey
	componentKey := commands[0]
	maps[cons.COMPONENT_KEY] = componentKey
	//第2个到最后一个参数为 component入参，少于2个参数肯定就是没入参,只有大于等于三个参数，才存在入参
	var params []string
	if len(commands) > 1 {
		params = commands[1:]
	}
	//入参解析
	for i := 0; i < len(params); i++ {
		str := params[i]
		//系统参数必定以--开头，后面跟若干字母
		if utils.RegularValidate(str, utils.REGEX_2DASE_AND_WORD) {
			if globalContext.findParameterType(componentKey, str) != cons.NO_VALUE {
				i++
				maps[str] = params[i]
			} else {
				maps[str] = componentKey
			}
			//组件参数必定以-开头，后面跟若干字母
		} else if utils.RegularValidate(str, utils.REGEX_DASE_AND_WORD) {
			p := str[1:]
			//只有单一参数才能跟值
			if len(p) == 1 {
				if globalContext.findParameterType(componentKey, str) != cons.NO_VALUE {
					i++
					maps[str] = params[i]
				} else {
					maps[str] = ""
				}
			} else {
				for m := 0; m < len(p); m++ {
					maps["-"+p[m:m+1]] = ""
				}
			}
		} else {
			return nil, errors.New("参数错误,参数应以-或--开头")
		}
	}
	//ENC(data)内容解密
	password, ok := maps[cons.SALT]
	if !ok && globalContext.Config != nil {
		password = globalContext.Config.GetString(cons.CONFIG_SALT)
	}
	if password != "" {
		for key, value := range maps {
			if value != "" && strings.HasPrefix(value, "ENC(") && strings.HasSuffix(value, ")") {
				value = value[4 : len(value)-1]
				value, err := utils.JasyptDec(value, password)
				if err != nil {
					return nil, errors.New("参数解密失败:" + err.Error())
				}
				maps[key] = value
			}
		}
	}
	return maps, nil
}

/*
*
添加config配置文件中的内容到入参中
*/
func addConfigToMap(maps map[string]string) {
	if globalContext.Config == nil {
		return
	}
	s, ok := maps[cons.COMPONENT_KEY]
	if !ok {
		return
	}
	cm := globalContext.findComponent(s, true)
	if cm == nil {
		return
	}
	for _, v := range cm.params {
		if _, ok = maps[v.CommandName]; v.ParamType != cons.NO_VALUE && v.ConfigName != "" && !ok {
			val := globalContext.Config.GetString(v.ConfigName)
			if val != "" {
				maps[v.CommandName] = val
				if strings.HasPrefix(val, "ENC(") && strings.HasSuffix(val, ")") {
					salt, ok := maps[cons.SALT]
					if !ok && globalContext.Config != nil {
						salt = globalContext.Config.GetString(cons.CONFIG_SALT)
					}
					val = val[4 : len(val)-1]
					val, err := utils.JasyptDec(val, salt)
					if err != nil {
						log.Error("参数解密失败", err.Error())
					}
					maps[v.CommandName] = val

				} else {
					maps[v.CommandName] = val
				}

			}
		}
	}
}
