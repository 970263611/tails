package basic

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/utils"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"sort"
	"strings"
)

type Context struct {
	components map[string]*componentMeta // Component 组件区
	cache      map[string]map[string]any //Cache 缓存数据区
	Config     *viper.Viper              //Config 配置信息
}

var globalContext *Context

func (c *Context) GetConfig() *viper.Viper {
	return c.Config
}

/*
*
设置缓存
*/
func (c *Context) SetCache(key string, value interface{}) bool {
	tid := utils.GetGoroutineID()
	if tid == "" {
		return false
	}
	gmap, ok := c.cache[tid]
	if !ok {
		gmap = make(map[string]any)
		c.cache[tid] = gmap
	}
	gmap[key] = value
	return true
}

/*
*
获取缓存数据
*/
func (c *Context) GetCache(key string) (interface{}, bool) {
	tid := utils.GetGoroutineID()
	if tid == "" {
		return nil, false
	}
	gmap, ok := c.cache[tid]
	if !ok {
		return nil, false
	}
	value, ok := gmap[key]
	return value, ok
}

/*
*
获取缓存数据
*/
func (c *Context) DelCache() {
	tid := utils.GetGoroutineID()
	if tid == "" {
		return
	}
	delete(c.cache, tid)
}

/*
*
配置信息转结构体
*/
func (c *Context) Unmarshal(a any) error {
	if c.Config != nil {
		if err := c.Config.Unmarshal(a); err != nil {
			return err
		}
	}
	return nil
}

/*
*
根据参数名称和查询方类型，返回组件参数名称
*/
func (c *Context) FindComponent(key string, isSystem bool) iface.ComponentMeta {
	return c.findComponent(key, isSystem)
}

/*
*
根据参数名称和查询方类型，返回组件参数名称
*/
func (c *Context) findComponent(key string, isSystem bool) *componentMeta {
	cm, ok := c.components[key]
	if ok {
		if !isSystem && cm.ComponentType == cons.EXECUTE {
			return cm
		} else {
			return cm
		}
	} else {
		cm = c.assembleByName(key)
		if cm != nil {
			return cm
		}
	}
	return nil
}

/*
*
根据组件名称及命令行参数名称，查询到指定参数的配置信息
*/
func (c *Context) findParameter(componentKey string, commandName string) (*parameter, error) {
	for _, value := range systemParam {
		arrs := strings.Split(value.CommandName, ",")
		for _, arr := range arrs {
			if arr == commandName {
				return value, nil
			}
		}
	}
	if componentKey != "" {
		cm := c.findComponent(componentKey, true)
		if cm != nil {
			for _, value := range cm.params {
				if value.CommandName == commandName {
					return value, nil
				}
			}
		}
	}
	msg := fmt.Sprintf("参数 %s 不存在", commandName)
	return nil, errors.New(msg)
}

/*
*
Help组件，传入组件名称，生成组件的help信息并返回；如果是查询tails支持的所有组件，则传入base
*/
func (c *Context) findHelp(key string) string {
	var sb strings.Builder
	var lineBreak, blank = "\n", ""
	if key == "base" {
		c.assembleAll()
		sb.WriteString("命令格式:")
		sb.WriteString(lineBreak)
		blank = utils.GetBlankByNum(2, " ")
		arr := []string{
			"获取全部组件列表：'base --help'",
			"获取组件帮助信息：'组件名称 --help'",
			"调用组件命令格式：'组件名称 组件参数列表'",
			"指定配置文件路径：'组件名称 --path 配置文件全路径 组件参数列表'",
			"组件转发请求：   '组件名称 --f或--forword ip:port或域名'",
			"传输参数需要解密：'组件名称 --salt 密钥 -p ENC(加密内容)' , 密钥也可配置在yml配置文件中,key为jasypt.salt",
			"多配置项目选择：  '组件名称 --key aabb',如果组件配置是web_svc.port,则此时取web_svc.aabb.port，仅组件参数支持此配置方式",
		}
		for _, v := range arr {
			sb.WriteString(blank + v)
			sb.WriteString(lineBreak)
		}
		sb.WriteString("组件列表:")
		components := c.components
		cnames := []string{}
		for k, _ := range components {
			cnames = append(cnames, k)
		}
		sort.Strings(cnames)
		for _, name := range cnames {
			component := components[name]
			sb.WriteString(lineBreak)
			sb.WriteString(utils.GetBlankByNum(2, " ") + component.GetName())
			sb.WriteString(lineBreak)
			describe := strings.Split(component.GetDescribe(), "\n")
			for _, v := range describe {
				sb.WriteString(utils.GetBlankByNum(4, " ") + v)
				sb.WriteString(lineBreak)
			}
		}
	} else {
		cm := globalContext.findComponent(key, false)
		if cm == nil {
			sb.WriteString(fmt.Sprintf("组件 %v 不存在，'./root --help' 查看组件列表", key))
		} else {
			sb.WriteString(lineBreak)
			sb.WriteString(fmt.Sprintf("组件 %v 参数列表:", key))
			sb.WriteString(lineBreak)
			params := cm.params
			if params != nil {
				for _, value := range cm.params {
					sb.WriteString(lineBreak)
					blank = utils.GetBlankByNum(2, " ")
					var requiredStr, upOrlowStr, configName string
					if value.Required {
						requiredStr = "(必要)"
					} else {
						requiredStr = "(可选)"
					}
					if strings.ToUpper(value.CommandName) == value.CommandName {
						upOrlowStr = " (大写)"
					}
					if value.ConfigName != "" {
						configName = cm.GetName() + "." + value.ConfigName + " "
					}
					sb.WriteString(blank + value.CommandName + "" + upOrlowStr + " " + requiredStr + " " + value.Describe + " ")
					sb.WriteString(lineBreak)
					if configName != "" {
						sb.WriteString(utils.GetBlankByNum(6, " ") + "配置文件对应key：" + configName)
						sb.WriteString(lineBreak)
					}
				}
			}
		}
	}
	return sb.String()
}

func (c *Context) addComponentMeta(key string, componentMeta *componentMeta) bool {
	c.components[key] = componentMeta
	return true
}

/*
*
注册待初始化组件
*/
func (c *Context) assembleByName(key string) *componentMeta {
	for _, cp := range initComponentList {
		if cp.GetName() == key {
			cm := newComponentMeta()
			cm.Component = cp
			var icm iface.ComponentMeta = cm
			cp.Register(icm)
			globalContext.addComponentMeta(key, cm)
			return cm
		}
	}
	return nil
}

/*
*
注册待初始化组件
*/
func (c *Context) assembleAll() {
	for _, cp := range initComponentList {
		_, ok := c.components[cp.GetName()]
		if !ok {
			c.assembleByName(cp.GetName())
		}
	}
}
