package config_update

import (
	"encoding/json"
	"encoding/xml"
	"gopkg.in/yaml.v3"
	"strings"
)

// 验证 JSON 文件的合法性
func isValidJSON(data []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(data, &js) == nil
}

// 验证 XML 文件的合法性（这里不假设特定的根元素）
func isValidXML(data []byte) bool {
	// 使用 xml.Unmarshal 尝试解析，但不指定具体的结构体
	var x xml.Token
	return xml.Unmarshal(data, &x) == nil && strings.HasPrefix(string(x.(xml.StartElement).Name.Local), "")
}

// 验证 YAML 文件的合法性
func isValidYAML(data []byte) bool {
	var y map[string]interface{}
	return yaml.Unmarshal(data, &y) == nil
}

// 验证 TEXT 文件的合法性
func isValidTEXT(data []byte) bool {
	return true
}

// 验证 HTML 文件的合法性
func isValidHTML(data []byte) bool {

	return true

}

// 验证 Properties 文件的合法性（键值对格式）
func isValidProperties(data []byte) bool {
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue // 跳过空行和注释行
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return false // 键值对格式不正确
		}
	}
	return true
}
