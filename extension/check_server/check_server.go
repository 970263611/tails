package check_server

import (
	"basic"
	othertool "basic/tool/other"
	"fmt"
	"math"
	"net"
	"strconv"
	"time"
)

type CheckServer struct{}

func GetInstance() basic.Component {
	return &CheckServer{}
}

func (c *CheckServer) GetOrder() int {
	return math.MaxInt64
}

func (c *CheckServer) Register() *basic.ComponentMeta {
	meta := &basic.ComponentMeta{
		Key:       "check_server",
		Describe:  "",
		Component: c,
	}
	return meta
}

func (c *CheckServer) Do(commands []string) (resp []byte) {
	var host string
	var port = -1
	length := len(commands)
	for i := 0; i < length; i++ {
		switch commands[i] {
		case "-h":
			i++
			if i < length {
				host = commands[i]
				if !othertool.IsValidIP(host) {
					return []byte("ip is not valid")
				}
			}
		case "-p":
			i++
			if i < length {
				p, err := strconv.Atoi(commands[i])
				if err != nil || p > 65535 {
					return []byte("port is not valid")
				}
				port = p
			}
		default:
			continue
		}
	}
	if len(host) == 0 || port < 0 {
		return []byte("parameter error")
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 3*time.Second)
	if err != nil {
		return []byte("disconnected")
	}
	defer conn.Close()
	return []byte("connected")
}
