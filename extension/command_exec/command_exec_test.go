package command_exec

import (
	"fmt"
	"testing"
)

func TestDoHandler(t *testing.T) {
	params := map[string]any{
		"command": "dir",
	}
	resp := doHandler(params)
	fmt.Println(string(resp))
}
