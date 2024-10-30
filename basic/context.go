package basic

// Component 组件区
var components = map[Key]*ComponentMeta{}

// Cache 缓存数据区
var cache = map[Key]any{}

// Config 配置信息
var config = map[Key]any{}

func Register(k Key, c *ComponentMeta) {
	components[k] = c
}
