package check_server

import (
	"fmt"
	"net"
	"time"
)

const key = "check_server"

func Register() (string, func(req map[string]any) (resp []byte), map[string]any, string) {
	return key, doHandler, nil, ""
}

func doHandler(req map[string]any) (resp []byte) {
	host := req["host"] // 替换为你要检查的主机名或IP地址
	port := req["port"] // Telnet服务默认端口号
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 3*time.Second)
	if err != nil {
		return []byte("disconnected")
	}
	defer conn.Close()
	return []byte("connected")
}
