package main

import (
	"basic"
	_ "basic/onload"
	"basic/tool/log_config"
	"fmt"
	"os"
)

func init() {
	config := log_config.LogConfig{OutType: 1}
	config.Init()
}

func main() {
	args := os.Args
	//调用component，并打印
	bytes := basic.Servlet(args)
	fmt.Println(string(bytes))
}
