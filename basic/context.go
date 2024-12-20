package basic

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/utils"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
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

func (c *Context) GetResponseWriter() http.ResponseWriter {
	cache, ok := c.GetCache(cons.WEB_RESP)
	if ok {
		writer := cache.(http.ResponseWriter)
		return writer
	}
	return nil
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
