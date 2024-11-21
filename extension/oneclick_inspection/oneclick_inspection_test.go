package oneclick_inspection

import (
	"fmt"
	"testing"
)

func TestDoHandler(t *testing.T) {
	params := map[string]any{
		"host":     "localhost",
		"port":     9999,
		"username": "sysadmin",
		"password": "Grcf#jk9527",
	}
	instance := GetInstance()
	resp := instance.Do(params)
	if resp != nil {
		fmt.Println(string(resp))
	}
}
