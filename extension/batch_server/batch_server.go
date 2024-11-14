package batch_server

import (
	"basic"
	"basic/tool/net"
	othertool "basic/tool/other"
	"errors"
	"fmt"
)

type BatchServer struct{}

func GetInstance() *BatchServer {
	return &BatchServer{}
}

func (b *BatchServer) GetName() string {
	return "batch_server"
}

func (b *BatchServer) GetDescribe() string {
	return "batch job 服务"
}

func (b *BatchServer) Register(globalContext *basic.Context) *basic.ComponentMeta {
	command := &basic.ComponentMeta{
		Component: b,
	}
	command.AddParameters(basic.STRING, "-h", "batch.server.ip", "host", true, func(s string) error {
		if !othertool.CheckIp(s) {
			return errors.New("batch服务的主机ip不合法")
		}
		return nil
	}, "batch服务的主机ip地址")
	command.AddParameters(basic.INT, "-p", "batch.server.port", "port", true, func(s string) error {
		if !othertool.CheckPortByString(s) {
			return errors.New("batch服务的端口port不合法")
		}
		return nil
	}, "batch服务的端口")
	command.AddParameters(basic.STRING, "-u", "batch.server.username", "username", true, nil, "batch登录用户名")
	command.AddParameters(basic.STRING, "-w", "batch.server.password", "password", true, nil, "batch登录密码")
	command.AddParameters(basic.STRING, "-q", "", "jobName", false, nil, "任务名称,通过任务名称查询任务状态")
	command.AddParameters(basic.STRING, "-e", "", "execJobByName", false, nil, "任务名称,通过任务名称重新执行任务")
	return command
}

func (b *BatchServer) Start(globalContext *basic.Context) error {
	return nil
}

func (b *BatchServer) Do(params map[string]any) (resp []byte) {
	host := params["host"].(string)
	port := params["port"].(int)
	urlPrefix := fmt.Sprintf("%s://%s:%d", "http", host, port)
	jobName, ok := params["jobName"].(string)
	if ok {
		m := map[string]string{
			"name": jobName,
		}
		var queryByNameResp QueryByNameResp
		err := net.PostRespStruct(urlPrefix+"/queryByName", m, nil, &queryByNameResp)
		if err != nil {
			fmt.Println("发错查询任务接口失败：", err)
			return
		}
		if queryByNameResp.Code == "0000000000000000" {
			if queryByNameResp.Data == "" {
				return []byte(fmt.Sprintf("该任务名称%v 未查询到", jobName))
			}
			return []byte(fmt.Sprintf("%v任务状态为: %v ", jobName, queryByNameResp.Data))
		} else {
			return []byte("查询接口失败:" + queryByNameResp.Data)
		}
	}
	return
}
