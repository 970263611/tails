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
