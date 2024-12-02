package edb_retry

import (
	"fmt"
	"testing"
)

func TestDoHandler(t *testing.T) {
	params := map[string]any{
		"host":            "localhost",
		"port":            9999,
		"username":        "sysadmin",
		"password":        "Grcf#jk9527",
		"receiveFileName": "receiveFile",
	}
	instance := GetInstance(nil)
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}
func TestD(t *testing.T) {

}
