package basic

type ComponentFunc func(req map[string]any) (resp []byte)

type Key string

type Config map[string]any

type ComponentMeta struct {
	ComponentFunc
	Key
	Config
	describe string
}

type Component interface {
	Register() (string, func(req map[string]any) (resp []byte), map[string]any, string)
}
