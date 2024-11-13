package basic

type ComponentType int

const (
	EXECUTE ComponentType = iota
	SYSTEM
)

type Component interface {
	GetName() string
	GetDescribe() string
	Register(c *Context) *ComponentMeta
	Start(c *Context) error
	Do(param map[string]any) (resp []byte)
}

/*
*
组件
*/
type ComponentMeta struct {
	ComponentType
	Component
	Params []Parameter
}

var initComponentList = make([]Component, 0)

/*
*
注册待初始化组件
*/
func AddInitComponent(cp Component) {
	if cp.GetName() != "" && cp.GetDescribe() != "" {
		initComponentList = append(initComponentList, cp)
	}
}

/*
*
组件全部参数验证
*/
func (c ComponentMeta) check(commands map[string]string) (map[string]interface{}, error) {
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
组件添加参数
*/
func (c *ComponentMeta) AddParameters(paramType ParamType, commandName string, configName string, standardName string, required bool, method func(s string) error, describe string) {
	p := Parameter{
		ParamType:    paramType,
		CommandName:  commandName,
		ConfigName:   configName,
		StandardName: standardName,
		Required:     required,
		CheckMethod:  method,
		Describe:     describe,
	}
	c.Params = append(c.Params, p)
}
