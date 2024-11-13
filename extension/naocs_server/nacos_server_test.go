package naocs_server

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
		"enabled":     "true",
		"namespaceId": "test1112",
		"serviceName": "cdl",
		"serviceIp":   "172.20.10.2",
		"servicePort": "9090",
	}
	instance := GetInstance()
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}
