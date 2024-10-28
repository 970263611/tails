package check_server

import (
	"fmt"
	"net"
	"time"
)

func Register() func(req map[string]any) (resp []byte) {
	return doHandler
}

func doHandler(req map[string]any) (resp []byte) {
	host := req["host"] // 替换为你要检查的主机名或IP地址
	port := req["port"] // Telnet服务默认端口号
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), 3*time.Second)
	defer conn.Close()
	if err != nil {
		return []byte("disconnected")
	}
	return []byte("connected")
}
