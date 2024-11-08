package main

import (
	"basic"
	"basic/log_config"
	_ "basic/onload"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func init() {
	//初始化日志配置
	config := log_config.NewLogConfig()
	config.OutType = 0
	log_config.Init(config)
}

func main() {
	//args := os.Args
	args := []string{"main.go", "web_server", "--start", "-p", "17085"}
	//解析命令行为map
	maps, err := basic.CommandsToMap(args)
	if err != nil {
		msg := fmt.Sprintf("入参解析失败 ： %v", err)
		log.Error(msg)
		fmt.Println(msg)
	}
	//调用component，并打印
	bytes := basic.Servlet(maps, false)
	fmt.Println(string(bytes))
}
