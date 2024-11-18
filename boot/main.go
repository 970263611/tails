package main

import (
	"basic"
	"basic/log_config"
	_ "basic/onload"
	"errors"
	"fmt"
	"os"
)

/*
*
系统参数  --help 获取组件帮助

		--path 指定配置文件路径
	    --addr ip:port或域名 时进行请求转发
*/
func main() {
	args := os.Args[1:]
	err := loadConfig(args)
	if err != nil {
		fmt.Println("加载配置文件失败:", err)
		return
	}
	//日志初始化
	initLogConfig()
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

type LogConfig struct {
	Log struct {
		Filename   string `mapstructure:"filename"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxBackups int    `mapstructure:"max_backups"`
		MaxAge     int    `mapstructure:"max_age"`
		Compress   bool   `mapstructure:"compress"`
		Level      string `mapstructure:"level"`
		OutType    int    `mapstructure:"out_type"`
	} `mapstructure:"log"`
}

/*
*
日志初始化
*/
func initLogConfig() {
	logconfig := &LogConfig{}
	basic.Unmarshal(logconfig)
	//初始化日志配置
	config := log_config.NewLogConfig()
	if logconfig.Log.Filename != "" {
		config.Filename = logconfig.Log.Filename
	}
	if logconfig.Log.MaxSize != 0 {
		config.MaxSize = logconfig.Log.MaxSize
	}
	if logconfig.Log.MaxBackups != 0 {
		config.MaxBackups = logconfig.Log.MaxBackups
	}
	if logconfig.Log.MaxAge != 0 {
		config.MaxAge = logconfig.Log.MaxAge
	}
	if logconfig.Log.Compress != false {
		config.Compress = logconfig.Log.Compress
	}
	if logconfig.Log.OutType != 0 {
		config.OutType = logconfig.Log.OutType
	}
	log_config.Init(config)
}
