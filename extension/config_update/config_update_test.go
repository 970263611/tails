package config_update

import (
	"fmt"
	"testing"
)

func TestDoHandler(t *testing.T) {
	params := map[string]any{
		"host":        "localhost",
		"port":        8848,
		"username":    "nacos",
		"password":    "nacos",
		"namespaceId": "test1112",
		"dataId":      "config-cpyy",
		"group":       "zhh",
		"outPut":      "D:/",
	}
	instance := GetInstance(nil)
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}

func TestUpdate(t *testing.T) {
	params := map[string]any{
		"host":        "localhost",
		"port":        8848,
		"username":    "nacos",
		"password":    "nacos",
		"namespaceId": "test1112",
		"dataId":      "config-cpyy",
		"group":       "zhh",
		"upload":      "D:/config-cpyy",
	}
	instance := GetInstance(nil)
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}
