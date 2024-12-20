package iface

import (
	cons "basic/constants"
	"github.com/spf13/viper"
	"net/http"
)

type Context interface {
	Servlet(commands []string, isSystem bool) ([]byte, error)
	LoadSystemParams(commands []string) ([]string, error)
	GetConfig() *viper.Viper
	SetCache(key string, value interface{}) bool
	GetCache(key string) (interface{}, bool)
	DelCache()
	GetResponseWriter() http.ResponseWriter
	Unmarshal(a any) error
	LoadConfig() error
	FindComponent(key string, isSystem bool) ComponentMeta
	FindSystemParams(key string) (string, bool)
	DelSystemParams(key string)
	SetForwordAddr(addr string)
}

type ComponentMeta interface {
	SetComponentType(componentType cons.ComponentType)
	AddParameters(paramType cons.ParamTypeIface, commandName cons.CommandNameIface, configName string, standardName string, required bool, method func(s string) error, describe string)
}

type Component interface {
	GetName() string
	GetDescribe() string
	Register(ComponentMeta)
	Do(param map[string]any) []byte
}
