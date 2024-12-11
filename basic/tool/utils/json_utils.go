package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// FormatJson 格式化json字符串
func FormatJson(marshal []byte) (res []byte, err error) {
	var out bytes.Buffer
	err = json.Indent(&out, marshal, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("格式化字符串失败: %v", err)
	}
	return out.Bytes(), nil
}
