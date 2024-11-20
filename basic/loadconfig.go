package basic

import (
	cons "basic/constants"
	"basic/tool/utils"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

/*
*
这里设置配置文件默认值
*/
var defaultParams = map[string]any{}

/*
*
配置文件解析
*/
func (c *Context) LoadConfig(path string) error {
	if c.Config != nil {
		return nil
	}
	v := viper.New()
	if path != "" {
		for key, value := range defaultParams {
			v.SetDefault(key, value)
		}
		v.SetConfigFile(path)
		if err := v.ReadInConfig(); err != nil {
			msg := fmt.Sprintf("配置文件 %v 读取失败，请检查路径和格式是否正确，仅支持yml或yaml格式文件,错误信息 : %v", path, err)
			log.Error(msg)
			return errors.New(msg)
		}
		msg := fmt.Sprintf("配置文件 %v 读取成功", path)
		log.Info(msg)
	}
	c.Config = v
	return nil
}

/*
*
系统参数列表
*/
var systemParam = []parameter{
	parameter{
		ParamType:    cons.NO_VALUE,
		CommandName:  "--help",
		StandardName: "",
		Required:     false,
		Describe:     "查看系统或组件帮助",
	},
	parameter{
		ParamType:    cons.STRING,
		CommandName:  "--path",
		StandardName: "",
		Required:     false,
		CheckMethod: func(s string) error {
			if s == "" {
				return errors.New("配置文件路径为空")
			}
			if !utils.FileExists(s) {
				msg := "配置文件不存在"
				log.Error(msg)
				return errors.New(msg)
			}
			if !utils.IsFileReadable(s) {
				msg := "配置文件不可读"
				log.Error(msg)
				return errors.New(msg)
			}
			return nil
		},
		Describe: "帮助",
	},
	parameter{
		ParamType:    cons.NO_VALUE,
		CommandName:  "--force",
		StandardName: "",
		Required:     false,
		Describe:     "强制执行",
	},
}
