package basic

type Component interface {
	GetOrder() int
	Register(gc *Context) *ComponentMeta
	Do(commands []string) (resp []byte)
}

type ComponentMeta struct {
	Component
	Key      string
	Describe string
}

func Assemble(list []Component) {
	for _, v := range list {
		cm := v.Register(globalContext)
		globalContext.Components[cm.Key] = cm
	}
}
