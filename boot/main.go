package main

import (
	"basic/tool/log_config"
	"flag"
	"fmt"
)

func init() {
	config := log_config.LogConfig{OutType: 0}
	config.Init()
}

func main() {
	//获得命令行参数
	//定义列名,定义默认数值,定义注释
	name := flag.String("name", "", "姓氏名谁")
	age := flag.Int("age", 0, "年龄")
	rmb := flag.Float64("rmb", 0, "资产")
	alive := flag.Bool("alive", true, "还活着呢吗")

	//解析胡获取参数
	flag.Parse()
	//丢入参数的指针中
	fmt.Println(*name, *age, *rmb, *alive)
}
