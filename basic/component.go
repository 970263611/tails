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
	Params   []*Parameter
}

type ParamType uint8

const (
	NO_VALUE ParamType = iota
	INT
	STRING
)

type Parameter struct {
	ParamType
	CommandName  string
	StandardName string
	Required     bool
	CheckMethod  func(string) error
	MinValue     int
	MaxValue     int
	RegexStr     string
	Describe     string
}

func Assemble(list []Component) {
	for _, v := range list {
		cm := v.Register(globalContext)
		globalContext.Components[cm.Key] = cm
	}
}
