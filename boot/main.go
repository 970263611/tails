package main

import (
	"basic"
	"basic/log_config"
	_ "basic/onload"
	"fmt"
	"os"
)

func init() {
	//初始化日志配置
	config := log_config.NewLogConfig()
	config.OutType = 0
	log_config.Init(config)
}

func main() {
	args := os.Args
	//调用component，并打印
	bytes := basic.Servlet(args, false)
	fmt.Println(string(bytes))
}
