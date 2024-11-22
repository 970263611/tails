package monitor_check

import (
	"basic/tool/net"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

type resultSingle struct {
	Desc  string
	Value string
}

type result struct {
	A1  resultSingle //昨日交易成功率 tps
	A2  resultSingle //当日实时成功率 tps
	A3  resultSingle //外部系统交易成功率 tps
	A4  resultSingle //子系统交易成功率 tps
	A5  resultSingle //应用进程检查
	A6  resultSingle //网络检查
	A10 resultSingle //统一监控报警检查
}

type findResult struct {
	*result
	host             string
	port             int
	username         string
	password         string
	token            string
	urlPrefix        string
	c                chan int
	version          string
	identifyCode     string
	identifyCodeUuid string
}

func (f findResult) login() (string, error) {
	var loginResp LoginResp
	var mapparams map[string]string
	if f.version == "2" {
		mapparams = map[string]string{
			"userId":           f.username,
			"passWord":         f.password,
			"identifyCode":     f.identifyCode,
			"identifyCodeUuid": f.identifyCodeUuid,
		}
	} else {
		mapparams = map[string]string{
			"userId":   f.username,
			"passWord": f.password,
		}
	}
	err := net.PostRespStruct(f.urlPrefix+"/api/cert/actions/login",
		mapparams, nil, &loginResp)
	if err != nil {
		return "登录失败", err
	}
	entries := loginResp.Data
	return entries.Token, nil
}

// 获取登录页面验证码
func (f findResult) getIdentifyCode() (*GetIdentifyCode, error) {
	var GetIdentifyCode GetIdentifyCode
	header := http.Header{}
	header.Set("Content-Type", "application/josn;charset=UTF-8")
	err := net.GetRespStruct(f.urlPrefix+"/api/cert/actions/identifyCode", nil, header, &GetIdentifyCode)
	if err != nil {
		return nil, err
	}
	return &GetIdentifyCode, nil
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
	queryParams.Add("selectTime", previousDay())
	f.result.A1.Desc = "昨日交易成功率"
	A1Resp, err := findSuccessRate(f, queryParams)
	if err != nil {
		f.result.A1.Value = fmt.Sprintf("查询失败 ： %v", err)
	}
	if A1Resp.Code == 0 {
		if A1Resp.Data == nil {
			f.result.A1.Value = "查询结果为空"
		} else {
			f.result.A1.Value = fmt.Sprint(A1Resp.Data[0].TotalSucRate)
		}
	} else {
		f.result.A1.Value = fmt.Sprintf("查询失败 ： %v", A1Resp.Message)
	}
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
	queryParams.Add("selectTime", currentDay())
	f.result.A2.Desc = "当日实时成功率和TPS"
	A1Resp, err := findSuccessRate(f, queryParams)
	if err != nil {
		f.result.A2.Value = fmt.Sprintf("查询失败 ： %v", err)
	}
	if A1Resp.Code == 0 {
		if A1Resp.Data == nil {
			f.result.A2.Value = "查询结果为空"
		} else {
			f.result.A2.Value = fmt.Sprintf("%v 和 %v", A1Resp.Data[0].TotalSucRate, A1Resp.Data[0].TradeTps)
		}
	} else {
		f.result.A2.Value = fmt.Sprintf("查询失败 ： %v", A1Resp.Message)
	}
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
	queryParams.Add("selectTime", currentDay())
	f.result.A3.Desc = "外部系统交易成功率"
	A1Resp, err := findSuccessRate(f, queryParams)
	if err != nil {
		f.result.A3.Value = fmt.Sprintf("查询失败 ： %v", err)
	}
	if A1Resp.Code == 0 {
		if A1Resp.Data == nil {
			f.result.A3.Value = "查询结果为空"
		} else {
			f.result.A3.Value = fmt.Sprint(A1Resp.Data[0].TotalSucRate)
		}
	} else {
		f.result.A3.Value = fmt.Sprintf("查询失败 ： %v", A1Resp.Message)
	}
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
	queryParams.Add("selectTime", currentDay())
	f.result.A4.Desc = "子系统交易成功率"
	A1Resp, err := findSuccessRate(f, queryParams)
	if err != nil {
		f.result.A4.Value = fmt.Sprintf("查询失败 ： %v", err)
	}
	if A1Resp.Code == 0 {
		if A1Resp.Data == nil {
			f.result.A4.Value = "查询结果为空"
		} else {
			f.result.A4.Value = fmt.Sprint(A1Resp.Data[0].TotalSucRate)
		}
	} else {
		f.result.A4.Value = fmt.Sprintf("查询失败 ： %v", A1Resp.Message)
	}
}

func (f findResult) a5() {
	defer f.deferMethod()
	queryParams := url.Values{}
	queryParams.Add("pageIndex", "1")
	queryParams.Add("pageSize", "50")
	queryParams.Add("centerFlag", "")
	header := http.Header{}
	header.Set("X-ER-UAT", f.token)
	f.result.A5.Desc = "应用进程检查"
	// 用于接收响应的结构体实例
	var A5Resp A5Resp
	// 发送 GET 请求
	err := net.GetRespStruct(f.urlPrefix+"/monitor/registration/get-node-status-list", queryParams, header, &A5Resp)
	if err != nil {
		f.result.A5.Value = fmt.Sprintf("查询失败 ： %v", err)
	}
	if A5Resp.Code == 0 {
		if A5Resp.Data == nil {
			f.result.A5.Value = "查询结果为空"
		} else {
			for i := 0; i < len(A5Resp.Data); i++ {
				status := A5Resp.Data[i].AlarmStatus
				statusStr := ""
				if status {
					statusStr = "正常"
				} else {
					statusStr = "异常"
				}
				f.result.A5.Value += fmt.Sprintf("%v服务，%v，最小节点数%v，当前节点数%v，当前异常节点数%v，当前异常ip %v； ",
					A5Resp.Data[i].SubsystemName,
					statusStr,
					A5Resp.Data[i].TotalNum,
					A5Resp.Data[i].CurrentNum,
					A5Resp.Data[i].ErrorNum,
					A5Resp.Data[i].ErrorIp,
				)
			}
		}
	} else {
		f.result.A5.Value = fmt.Sprintf("查询失败 ： %v", A5Resp.Message)
	}
}

func (f findResult) a6() {
	defer f.deferMethod()
	queryParams := url.Values{}
	queryParams.Add("pageIndex", "1")
	queryParams.Add("pageSize", "50")
	queryParams.Add("centerFlag", "")
	header := http.Header{}
	header.Set("X-ER-UAT", f.token)
	f.result.A6.Desc = "网络检查"
	// 用于接收响应的结构体实例
	var A6Resp A6Resp
	// 发送 GET 请求
	err := net.GetRespStruct(f.urlPrefix+"/monitor/tbMonctlPortConf/listPortAssist", queryParams, header, &A6Resp)
	if err != nil {
		f.result.A6.Value = fmt.Sprintf("查询失败 ： %v", err)
	}
	if A6Resp.Code == 0 {
		if A6Resp.Data == nil {
			f.result.A6.Value = "查询结果为空"
		} else {

			for i := 0; i < len(A6Resp.Data); i++ {
				status := A6Resp.Data[i].PortState
				statusStr := ""
				if status {
					statusStr = "正常"
				} else {
					statusStr = "异常"
				}
				f.result.A6.Value += fmt.Sprintf("%v中心，%v端口，%v； ",
					A6Resp.Data[i].CenterFlag,
					A6Resp.Data[i].HostPort,
					statusStr,
				)
			}
		}
	} else {
		f.result.A6.Value = fmt.Sprintf("查询失败 ： %v", A6Resp.Message)
	}
}

func (f findResult) a10() {
	defer f.deferMethod()
	queryParams := url.Values{}
	queryParams.Add("pageIndex", "1")
	queryParams.Add("pageSize", "50")
	queryParams.Add("alarmStatus", "0")
	header := http.Header{}
	header.Set("X-ER-UAT", f.token)
	f.result.A10.Desc = "统一监控报警检查"
	// 用于接收响应的结构体实例
	var A10Resp A10Resp
	// 发送 GET 请求
	err := net.GetRespStruct(f.urlPrefix+"/monitor/monitorAlarmJnl/list", queryParams, header, &A10Resp)
	if err != nil {
		f.result.A10.Value = fmt.Sprintf("查询失败 ： %v", err)
	}
	if A10Resp.Code == 0 {
		if A10Resp.Data == nil {
			f.result.A10.Value = "查询结果为空"
		} else {
			for i := 0; i < len(A10Resp.Data); i++ {
				level := A10Resp.Data[i].Level
				levelStr := ""
				switch level {
				case 0:
					levelStr = "正常"
					break
				case 1:
					levelStr = "预警"
					break
				case 2:
					levelStr = "告警"
					break
				case 3:
					levelStr = "严重告警"
					break
				default:
					break
				}
				f.result.A10.Value += fmt.Sprintf("%v指标，%v开始，%v； ",
					A10Resp.Data[i].JobName,
					formatTime(A10Resp.Data[i].BeginTime),
					levelStr,
				)
			}
		}
	} else {
		f.result.A10.Value = fmt.Sprintf("查询失败 ： %v", A10Resp.Message)
	}
}

func findSuccessRate(f findResult, queryParams url.Values) (*A1Resp, error) {
	header := http.Header{}
	header.Set("X-ER-UAT", f.token)
	// 用于接收响应的结构体实例
	var A1Resp A1Resp
	// 发送 GET 请求
	err := net.GetRespStruct(f.urlPrefix+"/monitor/trade/statistic", queryParams, header, &A1Resp)
	return &A1Resp, err
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

func formatTime(t string) string {
	layout := "2006-01-02T15:04:05.000"
	time, err := time.Parse(layout, t)
	if err != nil {
		log.Error("时间解析出错:", err.Error())
		return ""
	}
	return time.Format("2006-01-02 15:04:05")
}
