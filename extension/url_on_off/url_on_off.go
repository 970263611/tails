package uri_switch

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/utils"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

const (
	QUERY int = 1
	EXEC  int = 2
)

type SuspendServer struct{}

func GetInstance(globalContext iface.Context) iface.Component {
	return &SuspendServer{}
}

func (s *SuspendServer) GetName() string {
	return "url_on_off"
}

func (s *SuspendServer) GetDescribe() string {
	return "gateway流控交易接口服务 \n例：url_on_off  -h 127.0.0.1 -p 9999 -u gateway -w gateway -x 885e492d-bec8-4268-a2fc-ce0162c066fc -y GATEWAY-GSCT-TEST -n cfzt-edb -e false -U \"/cfzt-edb/query\" -g cfzt-edb -G cfzt-edb -c 0"
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

	cm.AddParameters(cons.STRING, cons.LOWER_Q, "query", "query", false, nil, "查询,传入api组名 查询API管理信息和流控信息")
	cm.AddParameters(cons.STRING, cons.LOWER_E, "enabled", "enabled", false, func(s string) error {
		if !utils.CheckIsBooleanByString(s) {
			return errors.New("enabled不合法,必须是boolean类型")
		}
		return nil
	}, "API组加入/移除uri: false加入, true是移除")
	cm.AddParameters(cons.STRING, cons.LOWER_G, "apiName", "apiName", false, nil, "uri所属API组名")
	cm.AddParameters(cons.STRING, cons.UPPER_U, "uri", "uri", false, nil, "要封停/解停的uri")
	cm.AddParameters(cons.INT, cons.LOWER_C, "count", "count", false, nil, "单机阈值(线程数)")
}

func (s *SuspendServer) Do(params map[string]any) (resp []byte) {
	res := &findResult{
		host:     params["host"].(string),
		port:     params["port"].(int),
		username: params["username"].(string),
		password: params["password"].(string),
		systemId: params["systemId"].(string),
		code:     params["code"].(string),
	}
	sel, err := check(params, res)
	if err != nil {
		return []byte(err.Error())
	}

	// 获取网关ID
	res.urlPrefix = fmt.Sprintf("%s://%s:%d", "http", res.host, res.port)
	token, err := res.login()
	if err != nil {
		return []byte("登录失败: " + err.Error())
	}
	if token == "" {
		return []byte("获取token: " + token)
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

	switch sel {
	case QUERY:
		return query(*res)
	case EXEC:
		return exec(*res)
	default:
		return []byte("未识别查询或操作指令,执行完毕")
	}
}

func check(params map[string]any, f *findResult) (int, error) {
	query, q := params["query"].(string)
	if q {
		f.apiName = query
		return QUERY, nil
	}
	enabled, e := params["enabled"].(string)
	if !e {
		return 0, errors.New("请输入enabled参数")
	} else {
		apiName, a := params["apiName"].(string)
		uri, u := params["uri"].(string)
		count, c := params["count"].(int)
		if a && u && c {
			f.enabled = enabled
			f.apiName = apiName
			f.uri = uri
			f.count = count
			return EXEC, nil
		} else {
			return 0, errors.New("请输入执行参数 apiName ,uri, count")
		}
	}
}

func query(res findResult) []byte {
	apiQueryRespEntry, err := res.selectApiGroup()
	if err != nil {
		return []byte("获取API组失败: " + err.Error())
	}
	releaseStatus := map[int]any{
		0: "精确",
		1: "前缀",
		2: "正则",
	}
	apiMap := map[string]any{
		"API名称": apiQueryRespEntry.ApiName,
		"匹配模式":  releaseStatus[apiQueryRespEntry.ReleaseStatus],
		"匹配串":   apiQueryRespEntry.PredicateItems,
	}
	apiByte, err := json.Marshal(apiMap)
	if err != nil {
		return []byte("解码失败: " + err.Error())
	}
	liuKong, err := res.queryLiuKong()
	if err != nil {
		return []byte("查询流控失败: " + err.Error())
	}
	resourceMode := map[int]any{
		0: "路由",
		1: "API组",
	}
	grade := map[int]any{
		0: "线程数",
		1: "QPS",
	}
	liuKongMap := map[string]any{
		"资源类型": resourceMode[liuKong.ResourceMode],
		"资源名称": liuKong.Resource,
		"阈值类型": grade[liuKong.Grade],
		"单机阈值": liuKong.Count,
	}
	liuKongByte, err2 := json.Marshal(liuKongMap)
	if err2 != nil {
		return []byte("解码失败: " + err.Error())
	}

	//格式化
	formatApiByte, err := utils.FormatJson(apiByte)
	if err != nil {
		return []byte("格式化失败: " + err.Error())
	}
	formatLiuKongByte, err := utils.FormatJson(liuKongByte)
	if err != nil {
		return []byte("格式化失败: " + err.Error())
	}
	sprintf := fmt.Sprintf("API组信息为:\n%s\n 流控信息为:\n%s", string(formatApiByte), string(formatLiuKongByte))
	return []byte(sprintf)
}

func exec(res findResult) []byte {
	// 1.封停解停uri, 取决于是否在api组中。如果在api组就是封停,如果不在就是解停
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

	// 2.创建流控,设置线程数
	liuKong, err := res.queryLiuKong()
	if err != nil {
		return []byte("获取流控失败: " + err.Error())
	}
	if liuKong == nil {
		respEntry, err := res.createLiuKong()
		if err != nil {
			return []byte("创建流控失败: " + err.Error())
		}
		if respEntry.Message != "" {
			return []byte("创建流控失败" + respEntry.Message)
		}
		log.Info("创建流控成功!")
	} else {
		res.liuKong = *liuKong
		respEntry, err := res.updateLiuKong()
		if err != nil {
			return []byte("更新流控失败: " + err.Error())
		}
		if respEntry.Message != "" {
			return []byte("更新流控失败" + respEntry.Message)
		}
		log.Info("更新流控成功!")
	}
	// 3.返回
	if res.enabled == "true" {
		return []byte(fmt.Sprintf("%v 解停成功", res.uri))
	}
	return []byte(fmt.Sprintf("%v 封停成功", res.uri))
}
