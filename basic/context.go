package basic

import (
	othertool "basic/tool/other"
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

type Context struct {
	components map[string]*ComponentMeta // Component 组件区
	cache      map[string]any            //Cache 缓存数据区
	Config     *viper.Viper              //Config 配置信息
}

var globalContext = &Context{
	components: make(map[string]*ComponentMeta),
	cache:      make(map[string]any),
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
func (c *Context) Unmarshal(a any) error {
	if err := globalContext.Config.Unmarshal(a); err != nil {
		return err
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

func (c *Context) FindHelp(key string) string {
	linewordnum := 100
	var msg string
	if key == "base" {
		c.assembleAll()
		//./root 组件名 组件参数列表
		components := c.components
		msg = "命令格式:\r\n"
		msg += "  组件调用: ./root compoment_key -params\r\n"
		msg += "  组件帮助: ./root compoment_key --help\r\n"
		msg += "组件列表:"
		for _, component := range components {
			msg += "\r\n"
			msg += "  " + component.GetName() + "\r\n"
			chinese := othertool.SplitByChinese(component.GetDescribe(), linewordnum)
			for i := 0; i < len(chinese); i++ {
				if i == 0 {
					msg += "      " + chinese[i] + "\r\n"
				} else {
					msg += "    " + chinese[i] + "\r\n"
				}
			}
		}
	} else {
		components := globalContext.FindComponent(key, false)
		if components == nil {
			msg = fmt.Sprintf("组件 %v 不存在，'./root --help' 查看组件列表", key)
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
	return msg
}

/*
*
命令行转map
*/
func CommandsToMap(commands []string) (map[string]string, error) {
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
			if globalContext.FindParameterType(componentKey, str) != NO_VALUE {
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
				if globalContext.FindParameterType(componentKey, str) != NO_VALUE {
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
