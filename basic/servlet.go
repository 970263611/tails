package basic

import (
	cons "basic/constants"
	"basic/tool/net"
	"basic/tool/utils"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

/*
*
组件分发
*/
func (c *Context) Servlet(commands []string, isSystem bool) ([]byte, error) {
	setGID()
	//参数中携带 --f或--forward ip:port或域名 时进行请求转发
	resp, err := forward(commands)
	if resp != nil || err != nil {
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
	if value := globalContext.FindSystemParams(cons.SYSPARAM_HELP); value != "" {
		return help(value), nil
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
	log.Infof("组件执行开始:[%v],组件入参:%v,调用组件参数个数:[%d]", key, commands, len(maps)-1)
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
func forward(commands []string) ([]byte, error) {
	var addr, params string
	//判断是否需要转发，并拼接转发参数
	addr = globalContext.FindSystemParams(cons.SYSPARAM_FORWORD)
	if addr == "" {
		return nil, nil
	}
	//转发功能仅单次有效
	globalContext.DelSystemParams(cons.SYSPARAM_FORWORD)
	//校验转发地址是否合法
	if !utils.CheckAddr(addr) {
		msg := fmt.Sprintf("地址不合法")
		log.Error(msg)
		return nil, errors.New(msg)
	}
	params = strings.Join(commands, " ")
	log.Infof("请求转发，转发地址:[%s],转发参数:[%s]", addr, params)
	//插入全局业务跟踪号
	gid := globalContext.FindSystemParams(cons.SYSPARAM_GID)
	if gid != "" {
		params += " " + cons.SYSPARAM_GID + " "
		params += gid
	}
	//转发请求并返回结果
	uri := fmt.Sprintf("http://%s/do", addr)
	resp, err := net.Post(uri, map[string]string{
		"params":   params,
		"isSystem": "",
	}, nil)
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer body.Close()
	str, err := io.ReadAll(body)
	if err != nil {
		log.Error("Error reading response body:", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(str))
	}
	return str, nil
}

/*
*
设置全局ID
*/
func setGID() {
	gid := globalContext.FindSystemParams(cons.SYSPARAM_GID)
	if gid == "" {
		gid, _ = utils.GenerateUUID()
		globalContext.setSystemParams(cons.SYSPARAM_GID, gid)
	}
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
	if len(commands) == 0 {
		return nil, nil
	}
	//第1个参数是 组件componentkey
	componentKey := commands[0]
	//go run main.go componentKey 至少应该有两个参数
	maps := make(map[string]string)
	maps[cons.COMPONENT_KEY] = componentKey
	//第2个到最后一个参数为 component入参，少于2个参数肯定就是没入参,只有大于等于三个参数，才存在入参
	var params []string
	if len(commands) > 1 {
		params = commands[1:]
	}
	//入参解析
	for i := 0; i < len(params); i++ {
		str := params[i]
		if utils.RegularValidate(str, utils.REGEX_DASE_AND_WORD) {
			p := str[1:]
			//只有单一参数才能跟值
			if len(p) == 1 {
				pt, err := globalContext.findParameter(componentKey, str)
				if err != nil {
					return nil, err
				} else if pt.ParamType != cons.NO_VALUE {
					i++
					if i >= len(params) {
						msg := fmt.Sprintf("参数 %s 解析失败,该参数必须有值", str)
						err = errors.New(msg)
						return nil, err
					}
					maps[str] = utils.RemoveQuotes(params[i])
				} else {
					maps[str] = ""
				}
			} else {
				for m := 0; m < len(p); m++ {
					maps["-"+p[m:m+1]] = ""
				}
			}
		}
	}
	//ENC(data)内容解密
	salt := globalContext.FindSystemParams(cons.SYSPARAM_SALT)
	if salt == "" && globalContext.Config != nil {
		salt = globalContext.Config.GetString(cons.CONFIG_SALT)
	}
	if salt != "" {
		for key, value := range maps {
			if value != "" && strings.HasPrefix(value, "ENC(") && strings.HasSuffix(value, ")") {
				value = value[4 : len(value)-1]
				value, err := utils.JasyptDec(value, salt)
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
			key := globalContext.FindSystemParams(cons.SYSPARAM_KEY)
			if key != "" {
				key = cm.GetName() + "." + key + "." + v.ConfigName
			} else {
				key = cm.GetName() + "." + v.ConfigName
			}
			val := globalContext.Config.GetString(key)
			if val != "" {
				maps[v.CommandName] = val
				if strings.HasPrefix(val, "ENC(") && strings.HasSuffix(val, ")") {
					salt := globalContext.FindSystemParams(cons.SYSPARAM_SALT)
					if salt == "" && globalContext.Config != nil {
						salt = globalContext.Config.GetString(cons.CONFIG_SALT)
					}
					val = val[4 : len(val)-1]
					val2, err := utils.JasyptDec(val, salt)
					if err != nil {
						log.Error("参数解密失败", err.Error())
					}
					maps[v.CommandName] = val2
				}
			}
		}
	}
}
