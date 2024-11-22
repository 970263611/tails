package batch_retry

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/net"
	"basic/tool/utils"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type BatchServer struct{}

func GetInstance(globalContext iface.Context) *BatchServer {
	return &BatchServer{}
}

func (b *BatchServer) GetName() string {
	return "batch_retry"
}

func (b *BatchServer) GetDescribe() string {
	return "batch job 服务,用于通过传入的任务名称 查询任务状态与重新执行任务功能 \n例：batch_retry -h 127.0.0.1 -p 9999 -u batch -w batch -s \"jobName\" "
}

func (b *BatchServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, cons.LOWER_H, "ip", "host", true, func(s string) error {
		if !utils.CheckIp(s) {
			return errors.New("batch服务的主机ip不合法")
		}
		return nil
	}, "batch服务的主机ip地址")
	cm.AddParameters(cons.INT, cons.LOWER_P, "port", "port", true, func(s string) error {
		if !utils.CheckPortByString(s) {
			return errors.New("batch服务的端口port不合法")
		}
		return nil
	}, "batch服务的端口")
	cm.AddParameters(cons.STRING, cons.LOWER_U, "username", "username", false, nil, "batch登录用户名")
	cm.AddParameters(cons.STRING, cons.LOWER_W, "password", "password", false, nil, "batch登录密码")
	cm.AddParameters(cons.STRING, cons.LOWER_S, "queryJobName", "queryJobName", false, nil, "任务名称,通过任务名称查询任务状态")
	cm.AddParameters(cons.STRING, cons.LOWER_R, "execJobName", "execJobName", false, nil, "任务名称,通过任务名称重新执行任务")
}

func (b *BatchServer) Do(params map[string]any) (resp []byte) {
	host := params["host"].(string)
	port := params["port"].(int)
	urlPrefix := fmt.Sprintf("%s://%s:%d", "http", host, port)
	jobName, ok := params["queryJobName"].(string)
	if ok {
		m := map[string]string{
			"name": jobName,
		}
		var resp Resp
		err := net.PostRespStruct(urlPrefix+"/queryByName", m, nil, &resp)
		if err != nil {
			log.Error("发错查询任务接口失败：", err)
			return []byte("发错查询任务接口失败：" + err.Error())
		}
		if resp.Code == "0000000000000000" {
			if resp.Data == "" {
				log.Error(fmt.Sprintf("该任务名称%v 未查询到", jobName))
				return []byte(fmt.Sprintf("该任务名称%v 未查询到", jobName))
			}
			return []byte(fmt.Sprintf("%v任务状态为: %v ", jobName, resp.Data))
		} else {
			return []byte("查询接口失败:" + resp.Data)
		}
	} else {
		log.Info("未传入查询任务名称...")
	}
	return []byte("批量任务执行结束")
}
