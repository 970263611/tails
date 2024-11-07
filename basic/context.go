package basic

import (
	"github.com/spf13/viper"
)

type Context struct {
	Components map[string]*ComponentMeta // Component 组件区
	Cache      map[string]any            //Cache 缓存数据区
	Config     *viper.Viper              //Config 配置信息
}

var globalContext = &Context{
	Components: make(map[string]*ComponentMeta),
	Cache:      make(map[string]any),
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

func FindComponent(key string, isSystem bool) *ComponentMeta {
	c, ok := globalContext.Components[key]
	if ok {
		if !isSystem && c.ComponentType == EXECUTE {
			return c
		} else {
			return c
		}
	}
	return nil
}
