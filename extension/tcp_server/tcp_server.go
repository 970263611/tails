package tcp_server

import (
	"basic"
	"os/exec"
	"strconv"
	"strings"
)

type TcpServer struct{}

func (t *TcpServer) GetDescribe() string {
	return "tcp连接服务"
}

func GetInstance() *TcpServer {
	return &TcpServer{}
}
func (t *TcpServer) GetName() string {
	return "tcp_server"
}

func (t *TcpServer) Register(globalContext *basic.Context) *basic.ComponentMeta {
	p1 := basic.Parameter{
		ParamType:    basic.INT,
		CommandName:  "-p",
		StandardName: "tcp_port",
		Required:     false,
		Describe:     "统计已连接上的，某端口tcp连接数",
	}
	p2 := basic.Parameter{
		ParamType:    basic.NO_VALUE,
		CommandName:  "-e",
		StandardName: "tcp_established",
		Required:     false,
		Describe:     "统计已连接上的，状态为“established的tcp连接数",
	}
	p3 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-i",
		StandardName: "tcp_ip",
		Required:     false,
		Describe:     "统计已连接上的，某ip的tcp连接数",
	}

	return &basic.ComponentMeta{
		//todo
		//ComponentType: basic.EXECUTE,
		Params:    []basic.Parameter{p1, p2, p3},
		Component: t,
	}
}

func (t *TcpServer) Start(globalContext *basic.Context) error {
	return nil
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
