package suspend_server

import (
	"basic/tool/net"
	"bytes"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
)

type findResult struct {
	host      string
	port      int
	username  string
	password  string
	systemId  string
	code      string
	name      string
	token     string
	urlPrefix string
	c         chan int
}

func (f findResult) login() (string, error) {
	var loginResp string
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=UTF-8")
	loginResp, err := net.PostRespString(f.urlPrefix+"/api/users/login",
		map[string]string{
			"userId":   f.username,
			"password": f.password,
		}, header)
	if err != nil {
		return "", err
	}
	v := viper.New()
	// 使用viper来解析JSON字符串
	v.SetConfigType("json") // 设置配置文件类型为json
	if err := v.ReadConfig(bytes.NewBuffer([]byte(loginResp))); err != nil {
		return "", err
	}
	token := v.GetString("token")
	return token, nil
}

// 获取网关配置ID
func (f findResult) selectGetWayID() (string, error) {
	queryParams := url.Values{}
	queryParams.Add("systemId", f.systemId)
	queryParams.Add("code", f.code)
	header := http.Header{}
	header.Set("Authorization", f.token)
	// 用于接收响应的结构体实例
	var SelectGetWayIDResp []SelectGetWayIDRespEntry
	// 发送 GET 请求
	err := net.GetRespStruct(f.urlPrefix+"/api/gateways/query", queryParams, header, &SelectGetWayIDResp)
	if err != nil {
		return "", err
	}
	for i := 0; i < len(SelectGetWayIDResp); i++ {
		if f.code == SelectGetWayIDResp[i].Code {
			return SelectGetWayIDResp[i].Id, nil
		}
	}
	return "", nil
}

// 获取路由配置信息
func (f findResult) selectRule(gatewayId string) (*SelectRuleRespEntry, error) {
	queryParams := url.Values{}
	queryParams.Add("gatewayId", gatewayId)
	queryParams.Add("name", f.name)
	header := http.Header{}
	header.Set("Authorization", f.token)
	// 用于接收响应的结构体实例
	var SelectRuleResp []SelectRuleRespEntry
	var SelectRuleRespEntry SelectRuleRespEntry
	// 发送 GET 请求
	err := net.GetRespStruct(f.urlPrefix+"/api/routes/query", queryParams, header, &SelectRuleResp)
	if err != nil {
		return &SelectRuleRespEntry, err
	}
	for i := 0; i < len(SelectRuleResp); i++ {
		if f.name == SelectRuleResp[i].Name {
			return &SelectRuleResp[i], nil
		}
	}
	return &SelectRuleRespEntry, nil
}
