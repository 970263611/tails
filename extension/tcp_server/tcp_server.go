package tcp_server

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"os/exec"
	"strconv"
	"strings"
)

type TcpServer struct{}

func (t *TcpServer) GetDescribe() string {
	return "tcp连接服务，例：tcp_server"
}

func GetInstance(globalContext iface.Context) iface.Component {
	return &TcpServer{}
}
func (t *TcpServer) GetName() string {
	return "tcp_server"
}

func (t *TcpServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.INT, cons.LOWER_P, "", "tcp_port", false, nil, "统计已连接上的，某端口tcp连接数，例：tcp_server -p 8080")
	cm.AddParameters(cons.NO_VALUE, cons.LOWER_E, "", "tcp_established", false, nil, "统计已连接上的，状态为“established的tcp连接数，例：tcp_server -e")
	cm.AddParameters(cons.STRING, cons.LOWER_H, "", "tcp_ip", false, nil, "统计已连接上的，某ip的tcp连接数，例：tcp_server -i 127.0.0.1")
}

func (t *TcpServer) Do(params map[string]any) (resp []byte) {
	cmd := exec.Command("netstat", "-ant")
	output, err := cmd.Output()
	if err != nil {
		return []byte("查询失败" + err.Error())
	}
	lines := strings.Split(string(output), "\n")
	count := 0

	//统计已连接上的，某端口tcp连接数
	tcp_portStr, tcp_port_ok := params["tcp_port"].(int)
	if tcp_port_ok {
		if tcp_portStr != 0 {
			for _, line := range lines {
				if strings.Contains(line, "ESTABLISHED") && strings.Contains(line, ":"+strconv.Itoa(tcp_portStr)) {
					count++
				}
			}
			return []byte(strconv.Itoa(count))
		}
	}

	//统计已连接上的，状态为“established的tcp连接数
	_, tcp_established_ok := params["tcp_established"].(string)
	if tcp_established_ok {
		for _, line := range lines {
			if strings.Contains(line, "ESTABLISHED") {
				count++
			}
		}
		return []byte(strconv.Itoa(count))
	}

	//统计已连接上的，某iptcp连接数
	tcp_ip, tcp_ip_ok := params["tcp_ip"].(string)
	if tcp_ip_ok {
		if len(tcp_ip) > 0 {
			for _, line := range lines {
				if strings.Contains(line, "ESTABLISHED") && strings.Contains(line, tcp_ip) {
					count++
				}
			}
			return []byte(strconv.Itoa(count))
		}
	}

	//默认查询
	return output
}
