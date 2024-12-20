package server_on_off

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
		"namespace":   "test1112",
		"serviceName": "cdl",
		"serviceIp":   "172.20.10.2",
		"servicePort": "9090",
		"groupName":   "",
	}
	instance := GetInstance(nil)
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}

func TestDoQuery(t *testing.T) {
	params := map[string]any{
		"host":        "localhost",
		"port":        8848,
		"username":    "nacos",
		"password":    "nacos",
		"query":       "",
		"enabled":     "false",
		"namespace":   "test1112",
		"serviceName": "cdl",
		"serviceIp":   "172.20.10.2",
		"servicePort": "9090",
		"groupName":   "",
		"clusterName": "DEFAULT",
	}
	instance := GetInstance(nil)
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}
