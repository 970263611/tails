package check_server

import (
	"basic"
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	args := []string{"main.go", "check_server", "-h", "127.0.0.1", "-p", "15431"}
	bytes := basic.Servlet(args)
	fmt.Println(string(bytes))
}
