package check_server

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	args := []string{"-h", "127.0.0.1", "-p", "15431"}
	ins := GetInstance()
	bytes := ins.Do(args)
	fmt.Println(string(bytes))
}
