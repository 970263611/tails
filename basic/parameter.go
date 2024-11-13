package basic

import (
	"fmt"
	"strconv"
)

type ParamType int8

const (
	NO_VALUE ParamType = iota
	INT
	STRING
)

/*
*
入参规则
*/
type Parameter struct {
	ParamType                         //参数的值类型，支持无值，整形和字符串三种
	CommandName  string               //命令行中的名称,例如 -h -p 等等
	ConfigName   string               //配置文件中的名称，例如 app.server.port
	StandardName string               //方法中的变量名,方法中获取的该参数时使用的key名称
	Required     bool                 //是否必填，默认false
	CheckMethod  func(s string) error //校验方法
	Describe     string               //描述
}

/*
*
参数验证
*/
func (p *Parameter) do(param map[string]any, commands map[string]string) error {
	err := p.check(commands)
	if err != nil {
		return err
	}
	err = p.paramFormat(param, commands)
	if err != nil {
		return err
	}
	return nil
}

/*
*
校验参数合法性
*/
func (p *Parameter) check(commands map[string]string) error {
	s, ok := commands[p.CommandName]
	//必填参数不存在直接报错
	if !ok && p.Required {
		return fmt.Errorf(p.StandardName + " is required")
	}
	//参数如果存在CheckMethod，则执行
	if ok && p.CheckMethod != nil {
		err := p.CheckMethod(s)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
*

	格式参数的类型，如果有方法内名称则替换为标准名称，并将参数添加到param中
*/
func (p *Parameter) paramFormat(param map[string]any, commands map[string]string) error {
	s, ok := commands[p.CommandName]
	if !ok {
		return nil
	}
	var a any = s
	if p.ParamType == INT {
		i, err := strconv.Atoi(fmt.Sprintf("%v", s))
		if err != nil {
			return fmt.Errorf("param %v format error : %v", p.CommandName, err)
		}
		a = i
	}
	if p.StandardName != "" {
		param[p.StandardName] = a
	} else {
		param[p.CommandName] = a
	}
	return nil
}
