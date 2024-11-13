package basic

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

/*
*
配置文件解析
*/
func LoadConfig(path string) error {
	v := viper.New()
	for key, value := range defaultParams {
		v.SetDefault(key, value)
	}
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		msg := fmt.Sprintf("配置文件 %v 读取失败，请检查路径和格式是否正确，仅支持yml或yaml格式文件", path)
		log.Error(msg)
		return errors.New(msg)
	}
	globalContext.Config = v
	msg := fmt.Sprintf("配置文件 %v 读取成功", path)
	log.Info(msg)
	return errors.New(msg)
}
