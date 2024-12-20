package basic_check

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/utils"
	"encoding/json"
	"errors"
	"sync"
)

const (
	CPYY        string = "CPYY"
	cpyyLogPath string = "/appcpyy/project/logs"

	GSCF        string = "GSCF"
	gscfLogPath string = "/appcpyy/project/logs"

	GRCF        string = "GRCF"
	grcfLogPath string = "/appcpyy/project/logs"

	logFileName string = "app.log"

	SUCCESS string = "应用启动成功"
)

type BasicServer struct {
	context iface.Context
}

func GetInstance(globalContext iface.Context) iface.Component {
	return &BasicServer{
		context: globalContext,
	}
}

func (b *BasicServer) GetName() string {
	return "basic_check"
}

func (b *BasicServer) GetDescribe() string {
	return "基础组项目投产检查(只支持基础组项目) " +
		"\n例：检查cfzt-edb项目是否一体化成功, 配置好cfzt-edb集群节点,因为当前节点需要到其他节点检查进程和日志,所以其他节点需要启动tails web" +
		"\n例：basic_check -a 127.0.0.1:6379,127.0.0.2:6379 -H 127.0.0.1:17001,127.0.0.2:17002 -n cfzt-edb -g CPYY"
}

func (b *BasicServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, cons.LOWER_A, "addr", "addr", true,
		func(s string) error {
			if !utils.CheckAddr(s) {
				return errors.New("主机地址不合法")
			}
			return nil
		}, "cfzt-xxx 主机地址")
	cm.AddParameters(cons.STRING, cons.UPPER_H, "tailsAddr", "tailsAddr", true, func(s string) error {
		if !utils.CheckAddr(s) {
			return errors.New("tails web地址不合法")
		}
		return nil
	}, "cfzt-xxx tailsWeb地址")
	cm.AddParameters(cons.STRING, cons.LOWER_N, "projectName", "projectName", true, nil, "基础组项目名称")
	cm.AddParameters(cons.STRING, cons.LOWER_G, "groupName", "groupName", true, nil, "部署到哪个项目组 CPYY,GSCF,GRCF")
}

func (b *BasicServer) Do(params map[string]any) (resp []byte) {
	req := &request{
		addr:        params["addr"].(string),
		tailsAddr:   params["tailsAddr"].(string),
		groupName:   params["groupName"].(string),
		projectName: params["projectName"].(string),
		basicServer: b,
	}
	//检查groupName 并且获取log地址
	err := req.getLogPath()
	if err != nil {
		return []byte(err.Error())
	}
	//初始化返回体
	req.initResp()

	//并发执行
	var wg sync.WaitGroup
	wg.Add(3)
	req.wg = wg
	go req.getLogInfo()
	go req.getPortInfo()
	go req.getInterfaceInfo()

	//等待返回
	req.wg.Wait()
	basicRespByte, err := json.Marshal(req.resp)
	if err != nil {
		return []byte("解码失败: " + err.Error())
	}
	formatRespByte, err := utils.FormatJson(basicRespByte)
	if err != nil {
		return []byte("格式化失败: " + err.Error())
	}
	return formatRespByte
}
