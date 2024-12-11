package edb_retry

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/net"
	"basic/tool/utils"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type EdbServer struct{}

func GetInstance(globalContext iface.Context) iface.Component {
	return &EdbServer{}
}

func (b *EdbServer) GetName() string {
	return "edb_retry"
}

func (b *EdbServer) GetDescribe() string {
	return "edb重发组件,用于文件处理异常恢复,针对于发送方与接收方. 发送方发送失败重发,接收方接收失败重发 " +
		"\n例: 发送方失败重发 	edb_retry -h 127.0.0.1 -p 9999 -u edb -w edb -n LCYXYS" +
		"\n例: 接收方失败重发	edb_retry -h 127.0.0.1 -p 9999 -u edb -w edb -N COMBINATION_ZONE"
}

func (b *EdbServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, cons.LOWER_H, "ip", "host", true, func(s string) error {
		if !utils.CheckIp(s) {
			return errors.New("edb服务的主机ip不合法")
		}
		return nil
	}, "edb服务的主机ip地址")
	cm.AddParameters(cons.INT, cons.LOWER_P, "port", "port", true, func(s string) error {
		if !utils.CheckPortByString(s) {
			return errors.New("edb服务的端口port不合法")
		}
		return nil
	}, "edb服务的端口")
	cm.AddParameters(cons.STRING, cons.LOWER_U, "username", "username", false, nil, "edb登录用户名")
	cm.AddParameters(cons.STRING, cons.LOWER_W, "password", "password", false, nil, "edb登录密码")
	cm.AddParameters(cons.STRING, cons.LOWER_N, "sendName", "sendName", false, nil, "财富中台-edb文件报送-接口名称")
	cm.AddParameters(cons.STRING, cons.UPPER_N, "receiveName", "receiveName", false, nil, "财富中台-edb文件接收-接口名称")
}

func (b *EdbServer) Do(params map[string]any) (resp []byte) {
	host := params["host"].(string)
	port := params["port"].(int)
	urlPrefix := fmt.Sprintf("%s://%s:%d", "http", host, port)
	sendFileName, sok := params["sendName"].(string)
	if sok {
		return sendfile(sendFileName, urlPrefix+"/edb/upload/retryByName")
	}
	receiveFileName, rok := params["receiveName"].(string)
	if rok {
		return sendfile(receiveFileName, urlPrefix+"/api/mainTask/reconsumeByName")
	}
	return
}

func sendfile(fileName, url string) (r []byte) {
	m := map[string]string{
		"dataName": fileName,
	}
	body := map[string]any{
		"body": m,
	}
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	var resp Resp
	err := net.PostRespStruct(url, body, header, &resp)
	if err != nil {
		log.Error("发错任务接口失败：", err)
		return []byte("发错任务接口失败：" + err.Error())
	}
	if resp.Code == "0000000000000000" {
		return []byte(fmt.Sprintf("%v文件重发成功: %v ", fileName, resp.Msg))
	} else {
		return []byte(fmt.Sprintf("%v文件重发失败: %v ", fileName, resp.Msg))
	}
}
