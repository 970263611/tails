package onload

import (
	"fmt"
	"os"
	"text/template"
)

type PluginData struct {
	Key string
	Pkg string
}

const pluginImportTemplate = `package onload
 
import (
  {{range .}}{{.Key}} "{{.Pkg}}"
  {{end}}
)
 
func init() {
  {{range .}}{{.Key}}.Register()
  {{end}}
}
`

// 插件列表
var plugins = []PluginData{
	PluginData{"pluginA", "day1022/pluginA"},
	PluginData{"pluginB", "day1022/pluginB"},
}

func Load() {
	// 创建模板并执行
	tmpl, err := template.New("importTemplate").Parse(pluginImportTemplate)
	if err != nil {
		fmt.Println("模板解析错误:", err)
		os.Exit(1)
	}

	// 创建文件
	file, err := os.Create("onload/onload.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, plugins); err != nil {
		fmt.Println("模板执行错误:", err)
		os.Exit(1)
	}
}
