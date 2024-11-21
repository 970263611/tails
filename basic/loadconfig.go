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
func (c *Context) LoadConfig() error {
	if c.Config != nil {
		return nil
	}
	v := viper.New()
	path := c.findSystemParams("--path")
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
var systemParam = []*parameter{
	&parameter{
		ParamType:    cons.NO_VALUE,
		CommandName:  cons.SYSPARAM_HELP,
		StandardName: "",
		Required:     false,
		Describe:     "查看系统或组件帮助",
	},
	&parameter{
		ParamType:    cons.STRING,
		CommandName:  cons.SYSPARAM_SALT,
		StandardName: "diffuse",
		Required:     false,
		Describe:     "参数key填充",
	},
	&parameter{
		ParamType:    cons.STRING,
		CommandName:  cons.SYSPARAM_PATH,
		StandardName: "",
		Required:     false,
		Describe:     "帮助",
	},
	&parameter{
		ParamType:    cons.STRING,
		CommandName:  cons.SYSPARAM_FORWORD,
		StandardName: "",
		Required:     false,
		Describe:     "查看系统或组件帮助",
	},
	&parameter{
		ParamType:    cons.STRING,
		CommandName:  cons.SYSPARAM_KEY,
		StandardName: "diffuse",
		Required:     false,
		Describe:     "查看系统或组件帮助",
	},
}

/*
*
解析入参中的系统参数 --开头的参数
*/
func (c *Context) LoadSystemParams(commands []string) ([]string, error) {
	var args []string
	var err error
	maps := make(map[string]any)
	if len(commands) == 0 {
		maps["--help"] = "base"
		args = commands
	} else if len(commands) == 1 {
		maps["--help"] = commands[0]
		args = commands
	} else {
		componemtKey := commands[0]
		if cm := c.findComponent(componemtKey, false); cm == nil {
			err = errors.New("组件 " + componemtKey + " 不存在")
		} else {
			args = append(args, componemtKey)
			params := commands[1:]
			for i := 0; i < len(params); i++ {
				str := params[i]
				if utils.RegularValidate(str, utils.REGEX_2DASE_AND_WORD) {
					pt, err2 := globalContext.findParameter("", str)
					if err2 != nil {
						err = err2
						args = commands
						break
					} else if pt.ParamType != cons.NO_VALUE {
						i++
						if i >= len(params) {
							msg := fmt.Sprintf("参数 %s 解析失败,该参数必须有值", str)
							err = errors.New(msg)
							break
						}
						maps[str] = params[i]
						if pt.StandardName == cons.DIFFUSE {
							args = append(args, str, params[i])
						}
					} else {
						maps[str] = componemtKey
					}
				} else {
					args = append(args, str)
				}
			}
		}
	}
	if err != nil {
		return args, err
	} else {
		c.SetCache(cons.SYSTEM_PARAMS, maps)
		return args, nil
	}
}

/*
*
查询入参当中的系统参数
*/
func (c *Context) findSystemParams(key string) string {
	cache, ok := c.GetCache(cons.SYSTEM_PARAMS)
	if ok {
		m := cache.(map[string]any)
		a, ok := m[key]
		if ok {
			s := a.(string)
			return s
		} else {
			return ""
		}
	} else {
		return ""
	}
}
