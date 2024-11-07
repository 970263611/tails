package basic

type ComponentType int

const (
	EXECUTE ComponentType = iota
	SYSTEM
)

type Component interface {
	GetOrder() int
	Register(c *Context) *ComponentMeta
	Do(param map[string]any) (resp []byte)
}

/*
*
组件
*/
type ComponentMeta struct {
	ComponentType
	Component
	Key      string
	Describe string
	Params   []Parameter
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

func (c *ComponentMeta) AddParameters(paramType ParamType, commandName string, standardName string, required bool, method func(s string) error, describe string) {
	p := Parameter{
		ParamType:    paramType,
		CommandName:  commandName,
		StandardName: standardName,
		Required:     required,
		CheckMethod:  method,
		Describe:     describe,
	}
	c.Params = append(c.Params, p)
}
