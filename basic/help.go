package basic

import (
	"basic/tool/utils"
	"fmt"
	"sort"
	"strings"
)

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
			"tails组件使用方法为'启动命令 组件名称 参数列表'，在windows下启动命令为tails.exe,在linux下启动命令为./tails,以下以linux内容均已linux命令为例。",
			"获取全组件列表信息：./tails",
			"获取指定组件帮助信息: ./tails 组件名称",
			"调用指定组件: ./tails 组件名称 参数列表",
			"通过web调用: http://ip:port/do?params=组件名称 参数列表",
			"指定配置文件: ./tails 组件名称 参数列表 --config或--c 配置文件全路径或相对路径，相对路径以tails文件所在目录为根路径，配置文件仅支持yml格式，--path与参数列表无前后顺序要求",
			"指令转发: ./tails 组件名称 参数列表 --f或--forward ip:port，--f或--forward与参数列表无前后顺序要求",
			"ENC解密传参：./tails 组件名称 --salt 密钥 参数列表(需加密内容用ENC()包裹)'，密钥也可配置在yml配置文件中,key为jasypt.salt",
			"动态多配置项方式配置：./tails 组件名称 --key 分组字符，如果组件配置是web_svc.port,分组字符是aabb,则此时取yml中web_svc.aabb.port",
			"日志级别配置: ./tails 组件名称 参数列表 --loglevel info   级别是info",
			"日志打印控制台配置: ./tails 组件名称 参数列表  --logconsole 打印到控制台",
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
			sb.WriteString(fmt.Sprintf("组件 %v 不存在，'./root ' 查看组件列表", key))
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
