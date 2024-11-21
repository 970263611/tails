package main

import (
	"basic"
	iface "basic/interfaces"
	"basic/log_config"
	"fmt"
	"os"
)

/*
*
系统参数  --help 获取组件帮助

	        --key 配置文件参数拼接
			--path 指定配置文件路径
		    --addr ip:port或域名 时进行请求转发
			--salt 解密密钥，解密方式jasypt-1.9.3.jar
*/
func main() {
	args := os.Args[1:]
	//创建全局context
	var context iface.Context = &basic.Context{}
	basic.InitGlobalContext(context)
	defer context.DelCache()
	//初始化组件
	basic.InitComponent()
	//解析入参中的系统参数
	args, err := context.LoadSystemParams(args)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//解析配置文件
	err = context.LoadConfig()
	if err != nil {
		fmt.Println("加载配置文件失败:", err)
		return
	}
	//日志初始化
	initLogConfig(context)
	//调用组件
	bytes, err := context.Servlet(args, false)
	//结果打印
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(bytes))
	}
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
func initLogConfig(c iface.Context) {
	logconfig := &LogConfig{}
	c.Unmarshal(logconfig)
	//初始化日志配置
	config := log_config.NewLogConfig(c)
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
