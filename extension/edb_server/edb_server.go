package edb_server

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/net"
	"basic/tool/utils"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type EdbServer struct{}

func GetInstance(globalContext iface.Context) iface.Component {
	return &EdbServer{}
}

func (b *EdbServer) GetName() string {
	return "edb_server"
}

func (b *EdbServer) GetDescribe() string {
	return "edb服务,用于文件处理异常恢复,针对于发送方与接收方. 发送方发送失败重发,接收方接收失败重发 例: edb_server -h 127.0.0.1 -p 9999 -u edb -w edb -s \"edbFile\""
}

func (b *EdbServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, "-h", "edb.server.ip", "host", true, func(s string) error {
		if !utils.CheckIp(s) {
			return errors.New("edb服务的主机ip不合法")
		}
		return nil
	}, "edb服务的主机ip地址")
	cm.AddParameters(cons.INT, "-p", "edb.server.port", "port", true, func(s string) error {
		if !utils.CheckPortByString(s) {
			return errors.New("edb服务的端口port不合法")
		}
		return nil
	}, "edb服务的端口")
	cm.AddParameters(cons.STRING, "-u", "edb.server.username", "username", false, nil, "edb登录用户名")
	cm.AddParameters(cons.STRING, "-w", "edb.server.password", "password", false, nil, "edb登录密码")
	cm.AddParameters(cons.STRING, "-s", "", "sendFileName", false, nil, "发送方需要恢复的文件名称")
	cm.AddParameters(cons.STRING, "-r", "", "receiveFileName", false, nil, "接收方需要恢复的文件名称")
}

func (b *EdbServer) Start() error {
	return nil
}

func (b *EdbServer) Do(params map[string]any) (resp []byte) {
	host := params["host"].(string)
	port := params["port"].(int)
	urlPrefix := fmt.Sprintf("%s://%s:%d", "http", host, port)
	sendFileName, sok := params["sendFileName"].(string)
	if sok {
		return sendfile(sendFileName, urlPrefix+"/api/mainTask/reconsumeByName")
	}
	receiveFileName, rok := params["receiveFileName"].(string)
	if rok {
		return sendfile(receiveFileName, urlPrefix+"/edb/upload/retryByName")
	}
	return
}

func sendfile(fileName, url string) (r []byte) {
	m := map[string]string{
		"dataName": fileName,
	}
	var resp Resp
	err := net.PostRespStruct(url, m, nil, &resp)
	if err != nil {
		log.Error("发错任务接口失败：", err)
		return []byte("发错任务接口失败：" + err.Error())
	}
	if resp.Code == "0000000000000000" {
		return []byte(fmt.Sprintf("%v文件重发成功: %v ", fileName, resp.Message))
	} else {
		return []byte(fmt.Sprintf("%v文件重发失败: %v ", fileName, resp.Message))
	}
}
