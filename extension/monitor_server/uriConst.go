package monitor_server

import (
	"basic/tool/net"
	"net/http"
	"net/url"
	"time"
)

type result struct {
	a1 any
	a2 any
	a3 any
	a4 any
}

type findResult struct {
	*result
	host      string
	port      int
	username  string
	password  string
	token     string
	urlPrefix string
	c         chan int
}

func (f findResult) login() (string, error) {
	var loginResp LoginResp
	err := net.PostRespStruct(f.urlPrefix+"/api/cert/actions/login",
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
	queryParams := url.Values{}
	queryParams.Add("gradingConfigId", "IN")
	queryParams.Add("pageIndex", "1")
	queryParams.Add("pageSize", "50")
	queryParams.Add("centerFlag", "")
	queryParams.Add("orderMethodId", "0")
	queryParams.Add("selectTimeDimension", "2")
	//昨日时间
	//queryParams.Add("selectTime", previousDay())
	findSuccessRate(f, queryParams)

}

func (f findResult) a2() {
	defer f.deferMethod()
	queryParams := url.Values{}
	queryParams.Add("gradingConfigId", "IN")
	queryParams.Add("pageIndex", "1")
	queryParams.Add("pageSize", "50")
	queryParams.Add("centerFlag", "")
	queryParams.Add("orderMethodId", "0")
	queryParams.Add("selectTimeDimension", "2")
	//当日时间
	/*queryParams.Add("selectTime", "2024-11-07+00:00:00")*/
	findSuccessRate(f, queryParams)
}

func (f findResult) a3() {
	defer f.deferMethod()
	queryParams := url.Values{}
	queryParams.Add("gradingConfigId", "OUT")
	queryParams.Add("pageIndex", "1")
	queryParams.Add("pageSize", "50")
	queryParams.Add("centerFlag", "")
	queryParams.Add("orderMethodId", "0")
	queryParams.Add("selectTimeDimension", "2")
	//当日时间
	/*queryParams.Add("selectTime", "2024-11-07+00:00:00")*/
	findSuccessRate(f, queryParams)
}

func (f findResult) a4() {
	defer f.deferMethod()
	queryParams := url.Values{}
	queryParams.Add("gradingConfigId", "INNER")
	queryParams.Add("pageIndex", "1")
	queryParams.Add("pageSize", "50")
	queryParams.Add("centerFlag", "")
	queryParams.Add("orderMethodId", "0")
	queryParams.Add("selectTimeDimension", "2")
	//当日时间
	/*queryParams.Add("selectTime", "2024-11-07+00:00:00")*/
	findSuccessRate(f, queryParams)
}

func (f findResult) a5() {
	defer f.deferMethod()
	queryParams := url.Values{}
	queryParams.Add("pageIndex", "1")
	queryParams.Add("pageSize", "50")
	queryParams.Add("centerFlag", "")
	header := http.Header{}
	header.Set("Authorization", f.token)
	// 用于接收响应的结构体实例
	var A5Resp A5Resp
	// 发送 GET 请求
	err := net.GetRespStruct(f.urlPrefix+"/monitor/registration/get-node-status-list", queryParams, header, &A5Resp)
	if err != nil {
		r := f.result
		r.a1 = "查询失败"
	}
}

func (f findResult) a6() {
	defer f.deferMethod()
	queryParams := url.Values{}
	queryParams.Add("pageIndex", "1")
	queryParams.Add("pageSize", "50")
	queryParams.Add("centerFlag", "")
	header := http.Header{}
	header.Set("Authorization", f.token)
	// 用于接收响应的结构体实例
	var A6Resp A6Resp
	// 发送 GET 请求
	err := net.GetRespStruct(f.urlPrefix+"/monitor/tbMonctlPortConf/listPortAssist", queryParams, header, &A6Resp)
	if err != nil {
		r := f.result
		r.a1 = "查询失败"
	}
}

func (f findResult) a10() {
	defer f.deferMethod()
	queryParams := url.Values{}
	queryParams.Add("pageIndex", "1")
	queryParams.Add("pageSize", "50")
	queryParams.Add("alarmStatus", "0")
	header := http.Header{}
	header.Set("Authorization", f.token)
	// 用于接收响应的结构体实例
	var A10Resp A10Resp
	// 发送 GET 请求
	err := net.GetRespStruct(f.urlPrefix+"/monitor/monitorAlarmJnl/list", queryParams, header, &A10Resp)
	if err != nil {
		r := f.result
		r.a1 = "查询失败"
	}
}

func findSuccessRate(f findResult, queryParams url.Values) {
	header := http.Header{}
	header.Set("Authorization", f.token)
	// 用于接收响应的结构体实例
	var A1Resp A1Resp
	// 发送 GET 请求
	err := net.GetRespStruct(f.urlPrefix+"/monitor/trade/statistic", queryParams, header, &A1Resp)
	if err != nil {
		r := f.result
		r.a1 = "查询失败"
	}
}

func currentDay() string {
	currentDay := time.Now()
	return currentDay.Format("2006-01-02") + "+00:00:00"
}

func previousDay() string {
	currentTime := time.Now()
	previousDay := currentTime.Add(-24 * time.Hour)
	previousDayFm := previousDay.Format("2006-01-02")
	return previousDayFm + "+00:00:00"
}
