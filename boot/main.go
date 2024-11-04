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
	//go run main.go component_key 至少应该有两个参数
	if len(args) < 2 {
		fmt.Println("parameter error")
		return
	}
	//最后一个参数是 component_key
	key := args[len(args)-1]
	//第二个至倒数第二个参数为 component入参，少于3个参数肯定就是没入参,只有大于等于三个参数，才存在入参
	var params []string
	if len(args) > 2 {
		params = args[1 : len(args)-1]
	}
	//调用component，并打印
	bytes := basic.Servlet(key, params)
	fmt.Println(string(bytes))
}
