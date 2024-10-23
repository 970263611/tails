package main

import (
	"basic/tool/log_config"
	"basic/tool/net"
	"fmt"
)

func init() {
	config := log_config.LogConfig{OutType: 0}
	config.Init()
}

func main() {
	handlers := make(map[string]func(req map[string]any) (resp []byte))
	handlers["/test1"] = handler
	handlers["/test2"] = handler
	net.Web(8080, handlers)
}

func handler(req map[string]any) (resp []byte) {
	fmt.Println(req)
	return []byte("Hello World")
}
