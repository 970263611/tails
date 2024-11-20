package basic

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"fmt"
	"strconv"
)

var initComponentList = make([]iface.Component, 0)

/*
*
注册待初始化组件
*/
func addInitComponent(cp iface.Component) {
	if cp.GetName() != "" && cp.GetDescribe() != "" {
		initComponentList = append(initComponentList, cp)
	}
}

/*
*
组件
*/
type componentMeta struct {
	cons.ComponentType
	iface.Component
	params []parameter
}

func newComponentMeta() *componentMeta {
	return &componentMeta{ComponentType: cons.EXECUTE, params: []parameter{}}
}

func (c *componentMeta) SetComponentType(componentType cons.ComponentType) {
	c.ComponentType = componentType
}

/*
*
组件添加参数
*/
func (c *componentMeta) AddParameters(paramType cons.ParamTypeIface, commandName cons.CommandNameIface, configName string, standardName string, required bool, method func(s string) error, describe string) {
	p := parameter{
		ParamType:    paramType.(cons.ParamType),
		CommandName:  commandName.GetCommandName(),
		ConfigName:   configName,
		StandardName: standardName,
		Required:     required,
		CheckMethod:  method,
		Describe:     describe,
	}
	c.params = append(c.params, p)
}

/*
*
组件全部参数验证
*/
func (c *componentMeta) check(commands map[string]string) (map[string]interface{}, error) {
	params := make(map[string]any)
	parameters := c.params
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
入参规则
*/
type parameter struct {
	cons.ParamType                      //参数的值类型，支持无值，整形和字符串三种
	CommandName    string               //命令行中的名称,例如 -h -p 等等
	ConfigName     string               //配置文件中的名称，例如 app.server.port
	StandardName   string               //方法中的变量名,方法中获取的该参数时使用的key名称
	Required       bool                 //是否必填，默认false
	CheckMethod    func(s string) error //校验方法
	Describe       string               //描述
}

/*
*
参数验证
*/
func (p *parameter) do(param map[string]any, commands map[string]string) error {
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
func (p *parameter) check(commands map[string]string) error {
	s, ok := commands[p.CommandName]
	//必填参数不存在直接报错
	if !ok && p.Required {
		return fmt.Errorf(p.StandardName + "是必要的")
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
func (p *parameter) paramFormat(param map[string]any, commands map[string]string) error {
	s, ok := commands[p.CommandName]
	if !ok {
		return nil
	}
	var a any = s
	if p.ParamType == cons.INT {
		i, err := strconv.Atoi(fmt.Sprintf("%v", s))
		if err != nil {
			return fmt.Errorf("参数 %v 格式化失败 : %v", p.CommandName, err)
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
