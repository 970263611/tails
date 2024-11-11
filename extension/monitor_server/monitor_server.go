package monitor_server

import (
	"basic"
	net "basic/tool/net"
	othertool "basic/tool/other"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type MonitorServer struct {
	*basic.Context
}

func GetInstance() *MonitorServer {
	return &MonitorServer{}
}

func (w *MonitorServer) GetName() string {
	return "monitor_server"
}

func (r *MonitorServer) GetDescribe() string {
	return "monitor监控指标服务"
}

func (r *MonitorServer) Register(globalContext *basic.Context) *basic.ComponentMeta {
	command := &basic.ComponentMeta{
		Component: r,
	}
	command.AddParameters(basic.STRING, "-h", "host", true, func(s string) error {
		if !othertool.CheckIp(s) {
			return errors.New("监控服务ip不合法")
		}
		return nil
	}, "监控服务的主机地址")
	command.AddParameters(basic.INT, "-p", "port", true, func(s string) error {
		if !othertool.CheckPortByString(s) {
			return errors.New("监控服务port不合法")
		}
		return nil
	}, "监控服务的端口")
	command.AddParameters(basic.STRING, "-u", "username", true, nil, "监控服务的登录用户名")
	command.AddParameters(basic.STRING, "-w", "password", true, nil, "监控服务的登录密码")
	return command
}
func (r *MonitorServer) Start(globalContext *basic.Context) error {
	return nil
}

var urlPrefix string

func (r *MonitorServer) Do(params map[string]any) (resp []byte) {
	//1. 获取ip和port组成请求url前缀
	host := params["host"]
	port := params["port"]
	urlPrefix = fmt.Sprintf("%s://%s:%d", "http", host, port)
	//2. 获取username和password作为请求体,发送post请求,获取返回token
	username := params["username"].(string)
	password := params["password"].(string)
	token, err := Login(username, password)
	if err != nil {
		log.Println("login接口请求失败:" + err.Error())
		return nil
	}
	//3. 携带token去做查询(开启线程并发查询)
	statisticResp, err := Statistic(token)
	if err != nil {
		log.Println("statistic接口请求失败:" + err.Error())
		return
	}
	entries2 := statisticResp.Data
	for _, entry := range entries2 {
		// 打印每个结构体的字段
		fmt.Printf("CenterFlg: %s, TotalNum: %d\n", entry.CenterFlg, entry.TotalSucCount)
	}
	//4. 并发查询其他接口
	//5. 组装查询数据并返回
	return nil
}

func Login(username, password string) (string, error) {
	postData := LoginReq{
		userId:   username,
		password: password,
	}
	var loginResp LoginResp
	err := net.PostRespStruct(urlPrefix+"/api/cert/actions/login", postData, nil, &loginResp)
	if err != nil {
		return "", err
	}
	entries := loginResp.Data
	return entries[0].Token, nil
}

func Statistic(token string) (*StatisticResp, error) {
	queryParams := url.Values{}
	queryParams.Add("pageIndex", "1")
	queryParams.Add("pageSize", "50")

	header := http.Header{}
	header.Set("Authorization", token)
	// 用于接收响应的结构体实例
	var statisticResp StatisticResp
	// 发送 GET 请求
	err := net.GetRespStruct(urlPrefix+"/monitor/trade/statistic", queryParams, header, &statisticResp)
	if err != nil {
		return nil, err
	}
	return &statisticResp, nil
}
