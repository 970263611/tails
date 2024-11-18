package suspend_server

import (
	"basic"
	"basic/tool/utils"
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
		if !utils.CheckIp(s) {
			return errors.New("微服务管理平台地址不合法")
		}
		return nil
	}, "微服务管理平台地址的主机地址")
	command.AddParameters(basic.INT, "-p", "suspend.server.port", "port", true, func(s string) error {
		if !utils.CheckPortByString(s) {
			return errors.New("微服务管理平台地址port不合法")
		}
		return nil
	}, "微服务管理平台地址的端口")
	command.AddParameters(basic.STRING, "-u", "suspend.server.username", "username", true, nil, "微服务管理平台的登录用户名")
	command.AddParameters(basic.STRING, "-w", "suspend.server.password", "password", true, nil, "微服务管理平台的登录密码")
	command.AddParameters(basic.STRING, "-s", "suspend.server.systemId", "systemId", true, nil, "系统ID")
	command.AddParameters(basic.STRING, "-c", "suspend.server.code", "code", true, nil, "网关编码")
	command.AddParameters(basic.STRING, "-n", "suspend.server.name", "name", true, nil, "组件名称")
	command.AddParameters(basic.STRING, "-e", "suspend.server.enabled", "enabled", true, func(s string) error {
		if !utils.CheckIsBooleanByString(s) {
			return errors.New("enabled不合法,必须是boolean类型")
		}
		return nil
	}, "封停/解停uri: false 是封停, true是解停")
	command.AddParameters(basic.STRING, "-U", "suspend.server.uri", "uri", true, nil, "要封停/解停的uri")
	command.AddParameters(basic.STRING, "-a", "suspend.server.apiName", "apiName", true, nil, "uri所属API组名")
	return command
}

func (s *SuspendServer) Start(globalContext *basic.Context) error {
	return nil
}
func (s *SuspendServer) Do(params map[string]any) (resp []byte) {
	res := &findResult{
		host:      params["host"].(string),
		port:      params["port"].(int),
		username:  params["username"].(string),
		password:  params["password"].(string),
		systemId:  params["systemId"].(string),
		code:      params["code"].(string),
		name:      params["name"].(string),
		gatewayId: params["gatewayId"].(string),
		enabled:   params["enabled"].(string),
		uri:       params["uri"].(string),
		apiName:   params["apiName"].(string),
		c:         make(chan int),
	}
	// 1.获取网关ID
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
	if gatewayId == "" {
		return []byte("获取网关配置ID是: " + gatewayId)
	}
	res.gatewayId = gatewayId

	// 2.封停解停uri, 取决于是否在api组中。如果在api组就是封停,如果不在就是解停
	apiQueryRespEntry, err := res.selectApiGroup()
	if err != nil {
		return []byte("获取API组失败: " + err.Error())
	}
	if apiQueryRespEntry == nil {
		if res.enabled == "true" {
			return []byte(fmt.Sprintf("%v 未加入限流无需解封", res.uri))
		}
		apiRespEntry, err := res.createApiGroup()
		if err != nil {
			return []byte("创建API组失败: " + err.Error())
		}
		if apiRespEntry.Message != "" {
			return []byte("创建API组失败" + apiRespEntry.Message)
		}
	} else {
		apiRespEntry, err := res.updateApiGroup(apiQueryRespEntry)
		if err != nil {
			return []byte("更新API组失败: " + err.Error())
		}
		if apiRespEntry.Message != "" {
			return []byte("更新API组失败: " + apiRespEntry.Message)
		}
	}

	// 3.创建流控,将api组设置线程数为0,代表不接受任何流量
	liuKong, err := res.queryLiuKong()
	if err != nil {
		return []byte("获取流控失败: " + err.Error())
	}
	if liuKong == nil {
		apiRespEntry, err := res.createLiuKong()
		if err != nil {
			return []byte("创建流控失败: " + err.Error())
		}
		if apiRespEntry.Message != "" {
			return []byte("创建流控失败" + apiRespEntry.Message)
		}
	}
	// 4.返回
	if res.enabled == "true" {
		return []byte(fmt.Sprintf("%v 解停成功", res.uri))
	}
	return []byte(fmt.Sprintf("%v 封停成功", res.uri))
}
