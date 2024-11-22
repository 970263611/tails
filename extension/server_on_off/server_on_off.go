package server_on_off

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/net"
	"basic/tool/utils"
	"errors"
	"fmt"
	"net/url"
)

type NacosServer struct {
	iface.Context
}

func GetInstance(globalContext iface.Context) iface.Component {
	return &NacosServer{}
}

func (w *NacosServer) GetName() string {
	return "server_on_off"
}

func (r *NacosServer) GetDescribe() string {
	return "Nacos上下线,用于封停/解停相关服务 \n例：server_on_off -h 127.0.0.1 -p 8848 -u nacos -w nacos -n test111 -e false -N sysName -H 127.0.0.1 -P 9090"
}

func (r *NacosServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, cons.LOWER_H, "ip", "host", true, func(s string) error {
		if !utils.CheckIp(s) {
			return errors.New("nacos服务ip不合法")
		}
		return nil
	}, "nacos服务的主机地址")
	cm.AddParameters(cons.INT, cons.LOWER_P, "port", "port", true, func(s string) error {
		if !utils.CheckPortByString(s) {
			return errors.New("nacos服务port不合法")
		}
		return nil
	}, "nacos服务的端口")
	cm.AddParameters(cons.STRING, cons.LOWER_U, "username", "username", true, nil, "nacos登录用户名")
	cm.AddParameters(cons.STRING, cons.LOWER_W, "password", "password", true, nil, "nacos登录密码")
	cm.AddParameters(cons.STRING, cons.LOWER_N, "namespace", "namespace", true, nil, "要封停/解停系统所在命名空间")
	cm.AddParameters(cons.STRING, cons.LOWER_E, "enabled", "enabled", true, func(s string) error {
		if !utils.CheckIsBooleanByString(s) {
			return errors.New("enabled不合法,必须是boolean类型")
		}
		return nil
	}, "是否要封停/解停, true是解停服务,false是封停服务")
	cm.AddParameters(cons.STRING, cons.UPPER_N, "serviceName", "serviceName", true, nil, "要封停/解停系统服务名")
	cm.AddParameters(cons.STRING, cons.LOWER_G, "groupName", "groupName", false, nil, "系统服务所在组名,不传默认DEFAULT_GROUP")
	cm.AddParameters(cons.STRING, cons.UPPER_H, "serviceIp", "serviceIp", true, nil, "要封停/解停系统ip")
	cm.AddParameters(cons.STRING, cons.UPPER_P, "servicePort", "servicePort", true, nil, "要封停/解停系统port")
}

func (r *NacosServer) Do(params map[string]any) (resp []byte) {
	host := params["host"].(string)
	port := params["port"].(int)
	urlPrefix := fmt.Sprintf("%s://%s:%d", "http", host, port)
	queryParams := url.Values{}
	queryParams.Add("username", params["username"].(string))
	queryParams.Add("password", params["password"].(string))
	queryParams.Add("enabled", params["enabled"].(string))
	queryParams.Add("namespaceId", params["namespace"].(string))
	queryParams.Add("serviceName", params["serviceName"].(string))
	queryParams.Add("ip", params["serviceIp"].(string))
	queryParams.Add("port", params["servicePort"].(string))
	if params["groupName"] != nil {
		queryParams.Add("groupName", params["groupName"].(string))
	}
	respString, err := net.PutRespString(urlPrefix+"/nacos/v1/ns/instance", queryParams, nil, nil)
	if err != nil {
		return []byte("发送Nacos上下线接口失败:, " + err.Error())
	}
	return []byte(respString)
}
