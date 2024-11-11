package monitor_server

import (
	"fmt"
	"testing"
)

func TestDoHandler(t *testing.T) {
	params := map[string]any{
		"host":     "127.0.0.1",
		"port":     9999,
		"username": "zhh",
		"password": "zhh",
	}
	instance := GetInstance()
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}
