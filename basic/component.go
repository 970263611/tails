package basic

import (
	"fmt"
	"strconv"
)

/*
*
组件
*/
type ComponentMeta struct {
	Component
	Key      string
	Describe string
	Params   []Parameter
}

type Component interface {
	GetOrder() int
	Register(gc *Context) *ComponentMeta
	Do(param map[string]any) (resp []byte)
}

type ParamType uint8

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
	ParamType
	CommandName  string               //命令行中的名称
	StandardName string               //方法中的变量名
	Required     bool                 //是否必填
	CheckMethod  func(s string) error //校验方法
	Describe     string               //描述
}

func Assemble(list []Component) {
	for _, v := range list {
		cm := v.Register(globalContext)
		globalContext.Components[cm.Key] = cm
	}
}

/*
*
组件全部参数验证
*/
func (c ComponentMeta) Check(commands map[string]string) (map[string]interface{}, error) {
	params := make(map[string]any)
	parameters := c.Params
	for _, param := range parameters {
		err := param.do(params, commands)
		if err != nil {
			return nil, err
		}
	}
	return params, nil
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
