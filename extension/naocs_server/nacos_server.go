package naocs_server

import (
	"basic"
	"basic/tool/net"
	othertool "basic/tool/other"
	"errors"
	"fmt"
	"net/url"
)

type NacosServer struct {
	*basic.Context
}

func GetInstance() *NacosServer {
	return &NacosServer{}
}

func (w *NacosServer) GetName() string {
	return "nacos_server"
}

func (r *NacosServer) GetDescribe() string {
	return "nacos服务"
}

func (r *NacosServer) Register(globalContext *basic.Context) *basic.ComponentMeta {
	command := &basic.ComponentMeta{
		Component: r,
	}
	command.AddParameters(basic.STRING, "-h", "host", true, func(s string) error {
		if !othertool.CheckIp(s) {
			return errors.New("nacos服务ip不合法")
		}
		return nil
	}, "nacos服务的主机地址")
	command.AddParameters(basic.INT, "-p", "port", true, func(s string) error {
		if !othertool.CheckPortByString(s) {
			return errors.New("nacos服务port不合法")
		}
		return nil
	}, "nacos服务的端口")
	command.AddParameters(basic.STRING, "-u", "username", true, nil, "nacos登录用户名")
	command.AddParameters(basic.STRING, "-w", "password", true, nil, "nacos登录密码")
	command.AddParameters(basic.STRING, "-e", "enabled", true, nil, "是否要封停/解停")
	command.AddParameters(basic.STRING, "-n", "namespaceId", true, nil, "要封停/解停系统所在命名空间")
	command.AddParameters(basic.STRING, "-s", "serviceName", true, nil, "要封停/解停系统服务名")
	command.AddParameters(basic.STRING, "-sh", "serviceIp", true, nil, "要封停/解停系统ip")
	command.AddParameters(basic.INT, "-sp", "servicePort", true, nil, "要封停/解停系统port")
	return command
}
func (r *NacosServer) Start(globalContext *basic.Context) error {
	return nil
}

func (r *NacosServer) Do(params map[string]any) (resp []byte) {
	host := params["host"].(string)
	port := params["port"].(int)
	urlPrefix := fmt.Sprintf("%s://%s:%d", "http", host, port)
	queryParams := url.Values{}
	queryParams.Add("username", params["username"].(string))
	queryParams.Add("password", params["password"].(string))
	queryParams.Add("enabled", params["enabled"].(string))
	queryParams.Add("namespaceId", params["namespaceId"].(string))
	queryParams.Add("serviceName", params["serviceName"].(string))
	queryParams.Add("ip", params["serviceIp"].(string))
	queryParams.Add("port", params["servicePort"].(string))

	respString, err := net.PutRespString(urlPrefix+"/nacos/v1/ns/instance", queryParams, nil)
	if err != nil {
		return []byte("发送Nacos上下线接口失败:, " + err.Error())
	}
	return []byte(respString)
}
