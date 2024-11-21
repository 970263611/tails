package server_check

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/utils"
	"errors"
	"fmt"
	"net"
	"time"
)

type CheckServer struct{}

func GetInstance(globalContext iface.Context) iface.Component {
	return &CheckServer{}
}

func (c *CheckServer) GetName() string {
	return "server_check"
}

func (c *CheckServer) GetDescribe() string {
	return "通过IP地址和端口号，判断应用服务是否正常或网络是否连通，例:svc_check -h 192.168.20.11 -p 15431"
}

func (c *CheckServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, cons.LOWER_H, "", "host", true,
		func(s string) error {
			if !utils.CheckIp(s) {
				return errors.New("IP不合法")
			}
			return nil
		}, "IP地址,支持ipv4,例:192.168.0.1")
	cm.AddParameters(cons.INT, cons.LOWER_P, "", "port", true,
		func(s string) error {
			if !utils.CheckPortByString(s) {
				return errors.New("端口不合法")
			}
			return nil
		}, "端口号,0 ~ 65535之间")
}

func (c *CheckServer) Do(params map[string]any) (resp []byte) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", params["host"].(string), params["port"].(int)), 3*time.Second)
	if err != nil {
		return []byte("服务不存在或网络不通")
	}
	defer conn.Close()
	return []byte("网络连通正常")
}
