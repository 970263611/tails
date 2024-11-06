package main

import (
	"basic"
	"basic/log_config"
	_ "basic/onload"
	"fmt"
)

func init() {
	//初始化日志配置
	config := log_config.NewLogConfig()
	config.OutType = 0
	log_config.Init(config)
}

func main() {
	//args := os.Args
	args := []string{"main.go", "check_server", "-h", "127.0.0.1", "-p", "154431", "--path", "config.yml"}
	//调用component，并打印
	bytes := basic.Servlet(args, false)
	fmt.Println(string(bytes))
}
