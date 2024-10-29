package check_server

import (
	"fmt"
	"net"
	"time"
)

const name = "check_server"

func Register() (key string, f func(req map[string]any) (resp []byte), metaData any) {
	return name, doHandler, nil
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
