package main

import (
	"basic"
	"basic/log_config"
	_ "basic/onload"
	"errors"
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
	err := loadConfig(args)
	if err != nil {
		fmt.Println("加载配置文件失败:", err)
		return
	}
	//调用component，并打印
	bytes := basic.Servlet(args, false)
	fmt.Println(string(bytes))
}

/*
*
加载配置文件,yaml
*/
func loadConfig(args []string) error {
	var path string
	for i := 0; i < len(args); i++ {
		if args[i] == "--path" {
			i++
			if i < len(args) {
				path = args[i]
			}
			if path == "" {
				return errors.New("传入配置文件路径为空")
			}
			break
		}
	}
	if path == "" {
		return nil
	}
	return basic.LoadConfig(path)
}
