package check_server

import (
	"fmt"
	"testing"
)

func TestDoHandler(t *testing.T) {
	params := map[string]any{
		"host": "127.0.0.1",
		"port": 11112,
	}
	resp := doHandler(params)
	fmt.Println(string(resp))
}
