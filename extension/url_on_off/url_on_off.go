package uri_switch

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/utils"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type SuspendServer struct{}

func GetInstance(globalContext iface.Context) iface.Component {
	return &SuspendServer{}
}

func (s *SuspendServer) GetName() string {
	return "url_on_off"
}

func (s *SuspendServer) GetDescribe() string {
	return "gateway封停解封交易接口服务 \n例：url_on_off  -h 127.0.0.1 -p 9999 -u gateway -w gateway -x 885e492d-bec8-4268-a2fc-ce0162c066fc -y GATEWAY-GSCT-TEST -n cfzt-edb -e false -U \"/cfzt-edb/query\" -g cfzt-edb "
}

func (s *SuspendServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, cons.LOWER_H, "ip", "host", true, func(s string) error {
		if !utils.CheckIp(s) {
			return errors.New("微服务管理平台地址不合法")
		}
		return nil
	}, "微服务管理平台地址的主机地址")
	cm.AddParameters(cons.INT, cons.LOWER_P, "port", "port", true, func(s string) error {
		if !utils.CheckPortByString(s) {
			return errors.New("微服务管理平台地址port不合法")
		}
		return nil
	}, "微服务管理平台地址的端口")
	cm.AddParameters(cons.STRING, cons.LOWER_U, "username", "username", true, nil, "微服务管理平台的登录用户名")
	cm.AddParameters(cons.STRING, cons.LOWER_W, "password", "password", true, nil, "微服务管理平台的登录密码")
	cm.AddParameters(cons.STRING, cons.LOWER_X, "systemId", "systemId", true, nil, "系统ID")
	cm.AddParameters(cons.STRING, cons.LOWER_Y, "code", "code", true, nil, "网关编码")
	cm.AddParameters(cons.STRING, cons.LOWER_N, "name", "name", true, nil, "组件名称")
	cm.AddParameters(cons.STRING, cons.LOWER_E, "enabled", "enabled", true, func(s string) error {
		if !utils.CheckIsBooleanByString(s) {
			return errors.New("enabled不合法,必须是boolean类型")
		}
		return nil
	}, "封停/解停uri: false 是封停, true是解停")
	cm.AddParameters(cons.STRING, cons.UPPER_U, "uri", "uri", true, nil, "(大写)要封停/解停的uri")
	cm.AddParameters(cons.STRING, cons.LOWER_G, "apiName", "apiName", true, nil, "uri所属API组名")
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
		enabled:  params["enabled"].(string),
		uri:      params["uri"].(string),
		apiName:  params["apiName"].(string),
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
		log.Info(fmt.Sprintf("创建%v API组成功!", res.apiName))
	} else {
		apiRespEntry, err := res.updateApiGroup(apiQueryRespEntry)
		if err != nil {
			return []byte("更新API组失败: " + err.Error())
		}
		if apiRespEntry.Message != "" {
			return []byte("更新API组失败: " + apiRespEntry.Message)
		}
		log.Info(fmt.Sprintf("更新%v API组成功!", res.apiName))
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
		log.Info("创建流控成功!")
	}
	// 4.返回
	if res.enabled == "true" {
		return []byte(fmt.Sprintf("%v 解停成功", res.uri))
	}
	return []byte(fmt.Sprintf("%v 封停成功", res.uri))
}
