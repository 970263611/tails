package monitor_server

import (
	"basic/tool/net"
	"net/http"
	"net/url"
)

type result struct {
	a1 any
	a2 any
	a3 any
	a4 any
}

type findResult struct {
	*result
	host     string
	port     int
	username string
	password string
	token    string
	c        chan int
}

func (f findResult) login() (string, error) {
	var loginResp LoginResp
	err := net.PostRespStruct(urlPrefix+"/api/cert/actions/login",
		map[string]string{
			"userId":   f.username,
			"password": f.password,
		}, nil, &loginResp)
	if err != nil {
		return "", err
	}
	entries := loginResp.Data
	return entries[0].Token, nil
}

func (f findResult) deferMethod() {
	f.c <- 0
}

func (f findResult) a1() {
	defer f.deferMethod()

}

func (f findResult) a2() {
	defer f.deferMethod()
}

func (f findResult) a3() {
	defer f.deferMethod()
}

func (f findResult) a4() {
	defer f.deferMethod()
}

func Statistic(token string) (string, error) {
	queryParams := url.Values{}
	queryParams.Add("gradingConfigId", "IN")
	queryParams.Add("pageIndex", "1")
	queryParams.Add("pageSize", "50")
	queryParams.Add("centerFlag", "")
	//queryParams.Add("orderMethodId", "0")
	//queryParams.Add("selectTimeDimension", "2")
	/*queryParams.Add("selectTime", "2024-11-07+00:00:00")*/

	header := http.Header{}
	header.Set("Authorization", token)
	// 发送 GET 请求
	resp, err := net.GetRespString(urlPrefix+"/monitor/trade/statistic", queryParams, header)
	if err != nil {
		return "", err
	}
	return resp, nil
}
