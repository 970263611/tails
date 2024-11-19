package iface

import (
	cons "basic/constants"
	"github.com/spf13/viper"
)

type Context interface {
	Servlet(commands []string, isSystem bool) []byte
	GetConfig() *viper.Viper
	SetCache(key string, value interface{}) bool
	GetCache(key string) (interface{}, bool)
	DelCache()
	Unmarshal(a any) error
	LoadConfig(path string) error
}

type ComponentMeta interface {
	SetComponentType(componentType cons.ComponentType)
	AddParameters(paramType cons.ParamType, commandName string, configName string, standardName string, required bool, method func(s string) error, describe string)
}

type Component interface {
	GetName() string
	GetDescribe() string
	Register(ComponentMeta)
	Start() error
	Do(param map[string]any) (resp []byte)
}
