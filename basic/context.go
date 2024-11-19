package basic

import (
	"basic/tool/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

type Context struct {
	components map[string]*ComponentMeta // Component 组件区
	cache      map[string]map[string]any //Cache 缓存数据区
	Config     *viper.Viper              //Config 配置信息
}

var globalContext = &Context{
	components: make(map[string]*ComponentMeta),
	cache:      make(map[string]map[string]any),
}

/*
*
设置缓存
*/
func SetCache(key string, value interface{}) bool {
	tid := utils.GetGoroutineID()
	if tid == "" {
		return false
	}
	gmap, ok := globalContext.cache[tid]
	if !ok {
		gmap = make(map[string]any)
		globalContext.cache[tid] = gmap
	}
	gmap[key] = value
	return true
}

/*
*
获取缓存数据
*/
func GetCache(key string) (interface{}, bool) {
	tid := utils.GetGoroutineID()
	if tid == "" {
		return nil, false
	}
	gmap, ok := globalContext.cache[tid]
	if !ok {
		return nil, false
	}
	return gmap[key], true
}

/*
*
获取缓存数据
*/
func DelCache() {
	tid := utils.GetGoroutineID()
	if tid == "" {
		return
	}
	delete(globalContext.cache, tid)
}

/*
*
这里设置配置文件默认值
*/
var defaultParams = map[string]any{}

/*
*
配置信息转结构体
*/
func Unmarshal(a any) error {
	if globalContext.Config != nil {
		if err := globalContext.Config.Unmarshal(a); err != nil {
			return err
		}
	}
	return nil
}

func (c *Context) addComponentMeta(key string, componentMeta *ComponentMeta) bool {
	c.components[key] = componentMeta
	return true
}

/*
*
注册待初始化组件
*/
func (c *Context) assembleByName(key string) *ComponentMeta {
	for _, cp := range initComponentList {
		if cp.GetName() == key {
			cm := cp.Register(globalContext)
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
			cm := cp.Register(globalContext)
			globalContext.addComponentMeta(cm.GetName(), cm)
		}
	}
}

/*
*
注册待初始化组件
*/
func (c *Context) Start() error {
	c.assembleAll()
	for _, cm := range c.components {
		err := cm.Start(c)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
*
根据参数名称和查询方类型，返回组件参数名称
*/
func (c *Context) FindComponent(key string, isSystem bool) *ComponentMeta {
	cm, ok := c.components[key]
	if ok {
		if !isSystem && cm.ComponentType == EXECUTE {
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
func (c *Context) FindParameterType(componentKey string, commandName string) ParamType {
	for _, value := range systemParam {
		if value.CommandName == commandName {
			return value.ParamType
		}
	}
	if componentKey != "" {
		componentMeta := c.FindComponent(componentKey, true)
		if componentMeta != nil {
			for _, value := range componentMeta.Params {
				if value.CommandName == commandName {
					return value.ParamType
				}
			}
		}
	}

	return -1
}

/*
*
Help组件，传入组件名称，生成组件的help信息并返回；如果是查询tails支持的所有组件，则传入base
*/
func (c *Context) FindHelp(key string) string {
	linewordnum := 100
	var sb strings.Builder
	var lineBreak string = "\r\n"
	if key == "base" {
		c.assembleAll()
		//./root 组件名 组件参数列表
		sb.WriteString("命令格式:")
		sb.WriteString(lineBreak)
		sb.WriteString(utils.GetBlankByNum(2, " ") + "获取全部组件列表: './boot base --help'")
		sb.WriteString(lineBreak)
		sb.WriteString(utils.GetBlankByNum(2, " ") + "获取组件帮助信息: './boot 组件名称 --help'")
		sb.WriteString(lineBreak)
		sb.WriteString(utils.GetBlankByNum(2, " ") + "调用组件命令格式: './boot 组件名称 组件参数列表'")
		sb.WriteString(lineBreak)
		sb.WriteString(utils.GetBlankByNum(2, " ") + "指定配置文件路径: './boot 组件名称 --path 配置文件全路径 组件参数列表'")
		sb.WriteString(lineBreak)
		sb.WriteString(utils.GetBlankByNum(2, " ") + "传输参数需要解密: './boot 组件名称 --salt 密钥 -p ENC(加密内容)' , 密钥也可配置在yml配置文件中,key为jasypt.salt")
		sb.WriteString(lineBreak)
		sb.WriteString(utils.GetBlankByNum(2, " ") + "通过web调用时忽略: './boot', 参数格式 params=组件名称 组件参数列表,支持post或get请求")
		sb.WriteString(lineBreak)
		sb.WriteString("组件列表")
		components := c.components
		for _, component := range components {
			sb.WriteString(lineBreak)
			sb.WriteString(utils.GetBlankByNum(2, " ") + component.GetName())
			sb.WriteString(lineBreak)
			chinese := utils.SplitByChinese(component.GetDescribe(), linewordnum)
			for i := 0; i < len(chinese); i++ {
				if i == 0 {
					sb.WriteString(utils.GetBlankByNum(6, " ") + chinese[i])
					sb.WriteString(lineBreak)
				} else {
					sb.WriteString(utils.GetBlankByNum(4, " ") + chinese[i])
					sb.WriteString(lineBreak)
				}
			}
		}
	} else {
		components := globalContext.FindComponent(key, false)
		if components == nil {
			sb.WriteString(fmt.Sprintf("组件 %v 不存在，'./root --help' 查看组件列表", key))
		} else {
			sb.WriteString("参数列表查看规则:")
			sb.WriteString(lineBreak)
			sb.WriteString(utils.GetBlankByNum(2, " ") + "'参数名 是否必须 配置文件中配置名(仅支持配置文件时才有)'，下方为参数描述")
			sb.WriteString(lineBreak)
			sb.WriteString(fmt.Sprintf("组件 %v 参数列表:", key))
			params := components.Params
			if params != nil {
				for _, value := range components.Params {
					sb.WriteString(lineBreak)
					if value.Required {
						sb.WriteString(utils.GetBlankByNum(2, " ") + value.CommandName + "  必要 " + value.ConfigName)
						sb.WriteString(lineBreak)
					} else {
						sb.WriteString(utils.GetBlankByNum(2, " ") + value.CommandName + "  可选 " + value.ConfigName)
						sb.WriteString(lineBreak)
					}
					sb.WriteString(utils.GetBlankByNum(6, " ") + value.Describe)
				}
			}
		}
	}
	return sb.String()
}

/*
*
添加config配置文件中的内容到入参中
*/
func addConfigToMap(maps map[string]string) {
	if globalContext.Config == nil {
		return
	}
	s, ok := maps[COMPONENT_KEY]
	if !ok {
		return
	}
	component := globalContext.FindComponent(s, true)
	if component == nil {
		return
	}
	for _, v := range component.Params {
		if _, ok = maps[v.CommandName]; v.ParamType != NO_VALUE && v.ConfigName != "" && !ok {
			val := globalContext.Config.GetString(v.ConfigName)
			if val != "" {
				if strings.HasPrefix(val, "ENC(") && strings.HasSuffix(val, ")") {
					salt, ok := maps[SALT]
					if !ok && globalContext.Config != nil {
						salt = globalContext.Config.GetString(CONFIG_SALT)
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
