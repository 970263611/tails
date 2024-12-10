package uri_switch

import (
	"basic/tool/net"
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"strconv"
)

type findResult struct {
	host      string
	port      int
	username  string
	password  string
	systemId  string
	code      string
	gatewayId string
	enabled   string
	apiName   string
	token     string
	urlPrefix string
	uri       string
	count     int
	liuKong   LiuKongQueryRespEntry
}

func (f findResult) login() (string, error) {
	var loginResp string
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=UTF-8")
	loginResp, err := net.PostRespString(f.urlPrefix+"/api/users/login",
		map[string]string{
			"username": f.username,
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

// 获取API组配置信息
func (f findResult) selectApiGroup() (*ApiQueryRespEntry, error) {
	queryParams := url.Values{}
	queryParams.Add("gatewayId", f.gatewayId)
	header := http.Header{}
	header.Set("Authorization", f.token)
	// 用于接收响应的结构体实例
	var apiQueryRespEntry []ApiQueryRespEntry
	// 发送 GET 请求
	err := net.GetRespStruct(f.urlPrefix+"/api/gateway-apidefinitions/query", queryParams, header, &apiQueryRespEntry)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(apiQueryRespEntry); i++ {
		if f.apiName == apiQueryRespEntry[i].ApiName {
			return &apiQueryRespEntry[i], nil
		}
	}
	return nil, nil
}

// 创建API组
func (f findResult) createApiGroup() (*RespEntry, error) {
	array := make([]any, 0)
	predicateItems := map[string]any{
		"pattern":       f.uri,
		"matchStrategy": 0,
	}
	array = append(array, predicateItems)
	jsonData, err := json.Marshal(array)
	if err != nil {
		log.Error("Error marshaling JSON:", err)
		return nil, err
	}
	postData := map[string]string{
		"predicateItems": string(jsonData),
		"apiName":        f.apiName,
		"gatewayId":      f.gatewayId,
	}

	header := http.Header{}
	header.Set("Authorization", f.token)
	header.Set("Content-Type", "application/json")
	var apiCreateRespEntry RespEntry
	// 发送 POST 请求
	err = net.PostRespStruct(f.urlPrefix+"/api/gateway-apidefinitions", postData, header, &apiCreateRespEntry)
	if err != nil {
		return nil, err
	}
	return &apiCreateRespEntry, nil
}

// 更新API组
func (f findResult) updateApiGroup(apiQuery *ApiQueryRespEntry) (*RespEntry, error) {
	items := apiQuery.PredicateItems
	// 1.将JSON字符串转换为map
	var data []map[string]interface{}
	err := json.Unmarshal([]byte(items), &data)
	if err != nil {
		log.Error("错误的解析了 PredicateItems JSON字符串:", err)
		return nil, err
	}
	// 2.找到要删除的元素的索引
	var indexToDelete = -1
	for i, item := range data {
		uri := item["pattern"]
		if uri == f.uri {
			indexToDelete = i
			break
		}
	}

	// 3.解停 删除data中的uri。封停 添加uri到data中
	if f.enabled == "true" {
		if indexToDelete == -1 {
			return &RespEntry{Message: fmt.Sprintf("%v 在%v组未查到,无法解停", f.uri, f.apiName)}, nil
		}
		data = append(data[:indexToDelete], data[indexToDelete+1:]...)
	} else {
		if indexToDelete == -1 {
			predicateItems := map[string]any{
				"pattern":       f.uri,
				"matchStrategy": 0,
			}
			data = append(data, predicateItems)
		}
	}
	// 4.将data转换回JSON字符串
	newJsonData, err := json.Marshal(data)
	if err != nil {
		log.Error("错误的解析了 JSON:", err)
		return nil, err
	}
	// 5.发送PUT请求,进行更新
	postData := map[string]string{
		"id":             apiQuery.Id,
		"predicateItems": string(newJsonData),
		"apiName":        apiQuery.ApiName,
		"releaseStatus":  strconv.Itoa(apiQuery.ReleaseStatus),
		"gatewayId":      apiQuery.GatewayId,
	}

	header := http.Header{}
	header.Set("Authorization", f.token)
	header.Set("Content-Type", "application/json")
	var apiUpdateRespEntry RespEntry
	err = net.PutRespStruct(f.urlPrefix+"/api/gateway-apidefinitions", nil, postData, header, &apiUpdateRespEntry)
	if err != nil {
		return nil, err
	}
	return &apiUpdateRespEntry, nil
}

// 查询流控
func (f findResult) queryLiuKong() (*LiuKongQueryRespEntry, error) {
	queryParams := url.Values{}
	queryParams.Add("gatewayId", f.gatewayId)
	header := http.Header{}
	header.Set("Authorization", f.token)
	// 用于接收响应的结构体实例
	var liuKongQueryRespEntry []LiuKongQueryRespEntry
	// 发送 GET 请求
	err := net.GetRespStruct(f.urlPrefix+"/api/gateway-flow-rules/query", queryParams, header, &liuKongQueryRespEntry)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(liuKongQueryRespEntry); i++ {
		if liuKongQueryRespEntry[i].ResourceMode == 1 && f.apiName == liuKongQueryRespEntry[i].Resource {
			return &liuKongQueryRespEntry[i], nil
		}
	}
	return nil, nil
}

// 创建流控
func (f findResult) createLiuKong() (*RespEntry, error) {
	liuKong := LiuKongQueryRespEntry{
		ResourceMode:         1,
		Resource:             f.apiName,
		ControlBehavior:      0,
		ParseStrategy:        0,
		FieldName:            "",
		Prop:                 false,
		Grade:                0,
		Count:                float32(f.count),
		IntervalSec:          1,
		Region:               "秒",
		Burst:                0,
		MatchStrategy:        0,
		RouteName:            f.apiName,
		MatchType:            "normal",
		Index:                0,
		MaxQueueingTimeoutMs: 1,
		GatewayId:            f.gatewayId,
	}
	header := http.Header{}
	header.Set("Authorization", f.token)
	header.Set("Content-Type", "application/json")
	var liuKongCreateRespEntry RespEntry
	// 发送 POST 请求
	err := net.PostRespStruct(f.urlPrefix+"/api/gateway-flow-rules", liuKong, header, &liuKongCreateRespEntry)
	if err != nil {
		return nil, err
	}
	return &liuKongCreateRespEntry, nil
}

func (f findResult) updateLiuKong() (*RespEntry, error) {
	liuKong := f.liuKong
	liuKong.Count = float32(f.count)
	header := http.Header{}
	header.Set("Authorization", f.token)
	header.Set("Content-Type", "application/json")
	var liuKongCreateRespEntry RespEntry
	// 发送 POST 请求
	err := net.PutRespStruct(f.urlPrefix+"/api/gateway-flow-rules", nil, liuKong, header, &liuKongCreateRespEntry)
	if err != nil {
		return nil, err
	}
	return &liuKongCreateRespEntry, nil
}
