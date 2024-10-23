package basic

type ComponentFunc func() (string, error)

type ComponentMeta struct {
	ComponentFunc
	Key      string
	MetaData any
}

type Component interface {
	Register() (key string, f func() (string, error), metaData any)
}
