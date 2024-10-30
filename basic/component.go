package basic

type Key string

type Describe string

type Component interface {
	GetOrder() int
	Register() *ComponentMeta
	Do(commands []string) (resp []byte)
}

type ComponentMeta struct {
	Key
	Describe
	Component
}
