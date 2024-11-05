package basic

type Component interface {
	GetOrder() int
	Register(gc *Context) *ComponentMeta
	Do(param map[string]any) (resp []byte)
}

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
