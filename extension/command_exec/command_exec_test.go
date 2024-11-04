package command_exec

import (
	"fmt"
	"testing"
)

func TestDoHandler(t *testing.T) {
	params := []string{"-c", "dir C:\\Windows"}
	ins := GetInstance()
	bytes := ins.Do(params)
	fmt.Println(string(bytes))
}
