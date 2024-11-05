package basic

type Context struct {
	Components map[string]*ComponentMeta // Component 组件区
	Cache      map[string]any            //Cache 缓存数据区
	Config     map[string]any            //Config 配置信息
}

var globalContext = &Context{
	Components: make(map[string]*ComponentMeta),
	Cache:      make(map[string]any),
	Config:     make(map[string]any),
}

/*
*
调用Register方法
*/
func Assemble(list []Component) {
	for _, v := range list {
		cm := v.Register(globalContext)
		globalContext.Components[cm.Key] = cm
	}
}

func FindParameterType(componentKey string, commandName string) ParamType {
	componentMeta, ok := globalContext.Components[componentKey]
	if ok {
		for _, value := range componentMeta.Params {
			if value.CommandName == commandName {
				return value.ParamType
			}
		}
	}
	return -1
}
