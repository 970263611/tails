package suspend_server

import (
	"basic"
	othertool "basic/tool/other"
	"encoding/json"
	"errors"
	"fmt"
)

type SuspendServer struct{}

func GetInstance() *SuspendServer {
	return &SuspendServer{}
}

func (s *SuspendServer) GetName() string {
	return "suspend_server"
}

func (s *SuspendServer) GetDescribe() string {
	return "封停解封交易服务"
}

func (s *SuspendServer) Register(globalContext *basic.Context) *basic.ComponentMeta {
	command := &basic.ComponentMeta{
		Component: s,
	}
	command.AddParameters(basic.STRING, "-h", "suspend.server.ip", "host", true, func(s string) error {
		if !othertool.CheckIp(s) {
			return errors.New("微服务管理平台地址不合法")
		}
		return nil
	}, "微服务管理平台地址的主机地址")
	command.AddParameters(basic.INT, "-p", "suspend.server.port", "port", true, func(s string) error {
		if !othertool.CheckPortByString(s) {
			return errors.New("微服务管理平台地址port不合法")
		}
		return nil
	}, "微服务管理平台地址的端口")
	command.AddParameters(basic.STRING, "-u", "suspend.server.username", "username", true, nil, "微服务管理平台的登录用户名")
	command.AddParameters(basic.STRING, "-w", "suspend.server.password", "password", true, nil, "微服务管理平台的登录密码")
	command.AddParameters(basic.STRING, "-s", "suspend.server.systemId", "systemId", true, nil, "系统ID")
	command.AddParameters(basic.STRING, "-c", "suspend.server.code", "code", true, nil, "网关编码")
	command.AddParameters(basic.STRING, "-n", "suspend.server.name", "name", true, nil, "组件名称")
	command.AddParameters(basic.STRING, "-U", "suspend.server.uri", "uri", true, nil, "要封停或解封得uri")
	return command
}

func (s *SuspendServer) Start(globalContext *basic.Context) error {
	return nil
}
func (s *SuspendServer) Do(params map[string]any) (resp []byte) {
	res := &findResult{
		host:     params["host"].(string),
		port:     params["port"].(int),
		username: params["username"].(string),
		password: params["password"].(string),
		systemId: params["systemId"].(string),
		code:     params["code"].(string),
		name:     params["name"].(string),
		c:        make(chan int),
	}
	res.urlPrefix = fmt.Sprintf("%s://%s:%d", "http", res.host, res.port)
	token, err := res.login()
	if err != nil {
		return []byte("登录失败: " + err.Error())
	}
	res.token = token
	gatewayId, err := res.selectGetWayID()
	if err != nil {
		return []byte("获取网关配置ID失败: " + err.Error())
	}
	SelectRuleRespEntry, err := res.selectRule(gatewayId)
	if err != nil {
		return []byte("获取路由配置信息: " + err.Error())
	}
	byteData, err := json.Marshal(SelectRuleRespEntry)
	if err != nil {
		fmt.Println("JSON序列化出错:", err)
		return
	}
	return byteData
}
