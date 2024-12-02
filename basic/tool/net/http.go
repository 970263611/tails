package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// GetRespString Get请求将resp封装到字符串/**
func GetRespString(urlStr string, params url.Values, headersMap http.Header) (string, error) {
	resp, err := Get(urlStr, params, headersMap)
	if err != nil {
		return "", err
	}
	body := resp.Body
	defer body.Close()
	str, err := io.ReadAll(body)
	if err != nil {
		log.Error("Error reading response body:", err)
		return "", err
	}
	return string(str), nil
}

// GetRespStruct Get请求将resp封装到结构体/**
func GetRespStruct(urlStr string, params url.Values, headersMap http.Header, response interface{}) error {
	resp, err := Get(urlStr, params, headersMap)
	if err != nil {
		return err
	}
	body := resp.Body
	defer body.Close()
	// 解码响应体到结构体
	err = json.NewDecoder(body).Decode(&response)
	if err != nil {
		return fmt.Errorf("解码响应体失败: %v", err)
	}
	return nil
}

// PostRespString Post请求将resp封装到字符串/**
func PostRespString(urlStr string, postData interface{}, headersMap http.Header) (string, error) {
	resp, err := Post(urlStr, postData, headersMap)
	if err != nil {
		return "", err
	}
	body := resp.Body
	defer body.Close()
	// 读取响应体
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %v", err)
	}
	return string(bodyBytes), nil
}

// PostRespStruct Post请求将resp封装到结构体/**
func PostRespStruct(urlStr string, postData interface{}, headersMap http.Header, response interface{}) error {
	resp, err := Post(urlStr, postData, headersMap)
	if err != nil {
		return err
	}
	body := resp.Body
	defer body.Close()
	// 解码响应体到结构体
	err = json.NewDecoder(body).Decode(&response)
	if err != nil {
		return fmt.Errorf("解码响应体失败: %v", err)
	}
	return nil
}

// PutRespString Put请求将resp封装到字符串/**
func PutRespString(urlStr string, params url.Values, postData interface{}, headersMap http.Header) (string, error) {
	resp, err := Put(urlStr, params, postData, headersMap)
	if err != nil {
		return "", err
	}
	body := resp.Body
	defer body.Close()
	str, err := io.ReadAll(body)
	if err != nil {
		log.Error("Error reading response body:", err)
		return "", err
	}
	return string(str), nil
}

// PutRespStruct Put请求将resp封装到结构体/**
func PutRespStruct(urlStr string, params url.Values, postData interface{}, headersMap http.Header, response interface{}) error {
	resp, err := Put(urlStr, params, postData, headersMap)
	if err != nil {
		return err
	}
	body := resp.Body
	defer body.Close()
	// 解码响应体到结构体
	err = json.NewDecoder(body).Decode(&response)
	if err != nil {
		return fmt.Errorf("解码响应体失败: %v", err)
	}
	return nil
}

// Get 通用Get请求,返回原生resp/**
func Get(urlStr string, params url.Values, headersMap http.Header) (*http.Response, error) {
	// 构建完整的 URL（如果有查询参数）
	var fullURL string
	if params != nil {
		u, err := url.Parse(urlStr)
		if err != nil {
			return nil, fmt.Errorf("解析 URL 失败: %v", err)
		}
		u.RawQuery = params.Encode()
		fullURL = u.String()
	} else {
		fullURL = urlStr
	}
	// 创建请求
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	// 遍历headersMap并添加到请求的Header中
	for key, values := range headersMap {
		joinedValues := strings.Join(values, ",")
		req.Header.Set(key, joinedValues)
	}
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	req, resp, err = GetLog(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Post 通用Post请求,返回原生resp/**
func Post(urlStr string, postData interface{}, headersMap http.Header) (*http.Response, error) {
	// 创建请求体
	var reqBody *strings.Reader
	if postData != nil {
		postDataBytes, err := json.Marshal(postData)
		if err != nil {
			return nil, fmt.Errorf("编码 POST 数据失败: %v", err)
		}
		reqBody = strings.NewReader(string(postDataBytes))
	} else {
		return nil, fmt.Errorf("post请求,请求体不能为null")
	}
	// 创建请求
	req, err := http.NewRequest(http.MethodPost, urlStr, reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// 遍历headersMap并添加到请求的Header中
	for key, values := range headersMap {
		joinedValues := strings.Join(values, ",")
		req.Header.Set(key, joinedValues)
	}
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	req, resp, err = PostLog(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Put 通用Put请求,返回原生resp/**
func Put(urlStr string, params url.Values, postData interface{}, headersMap http.Header) (*http.Response, error) {
	// 构建完整的 URL（如果有查询参数）
	var fullURL string
	if params != nil {
		u, err := url.Parse(urlStr)
		if err != nil {
			return nil, fmt.Errorf("解析 URL 失败: %v", err)
		}
		u.RawQuery = params.Encode()
		fullURL = u.String()
	} else {
		fullURL = urlStr
	}
	// 创建请求体
	var reqBody *strings.Reader
	if postData != nil {
		postDataBytes, err := json.Marshal(postData)
		if err != nil {
			return nil, fmt.Errorf("编码 PUT 数据失败: %v", err)
		}
		reqBody = strings.NewReader(string(postDataBytes))
	}
	var req *http.Request
	var err error
	if reqBody == nil {
		req, err = http.NewRequest(http.MethodPut, fullURL, nil)
	} else {
		req, err = http.NewRequest(http.MethodPut, fullURL, reqBody)
	}
	// 遍历headersMap并添加到请求的Header中
	for key, values := range headersMap {
		joinedValues := strings.Join(values, ",")
		req.Header.Set(key, joinedValues)
	}
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	req, resp, err = PutLog(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// post 打印请求报文和响应报文
func PostLog(req *http.Request, resp *http.Response) (*http.Request, *http.Response, error) {
	reqbodyBytes, err1 := ioutil.ReadAll(req.Body)
	if err1 != nil {
		log.Error("读取请求体出错:", err1)
		return req, resp, fmt.Errorf("读取请求体出错: %v", err1)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(reqbodyBytes))

	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Error("读取响应体出错:", err2)
		return req, resp, fmt.Errorf("读取响应体出错: %v", err2)
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	log.Info("请求方法: ", req.Method, "\n请求URL:", req.URL, "\n请求头:", req.Header, "\n请求体:", string(reqbodyBytes),
		"\n响应状态码:", resp.StatusCode, "\n响应头:", resp.Header, "\n响应体:", string(bodyBytes))
	return req, resp, nil
}

// get 打印请求报文和响应报文
func GetLog(req *http.Request, resp *http.Response) (*http.Request, *http.Response, error) {
	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Error("读取响应体出错:", err2)
		return req, resp, fmt.Errorf("读取响应体出错: %v", err2)
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	log.Info("请求方法: ", req.Method, "\n请求URL:", req.URL, "\n请求头:", req.Header,
		"\n响应状态码:", resp.StatusCode, "\n响应头:", resp.Header, "\n响应体:", string(bodyBytes))
	return req, resp, nil
}

// put 打印请求报文和响应报文
func PutLog(req *http.Request, resp *http.Response) (*http.Request, *http.Response, error) {
	reqbodyBytes := []byte("")
	if req.Body != nil {
		reqbodyBytesTemp, err1 := ioutil.ReadAll(req.Body)
		if err1 != nil {
			log.Error("读取请求体出错:", err1)
			return req, resp, fmt.Errorf("读取请求体出错: %v", err1)
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(reqbodyBytesTemp))
		reqbodyBytes = reqbodyBytesTemp
	}
	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Error("读取响应体出错:", err2)
		return req, resp, fmt.Errorf("读取响应体出错: %v", err2)
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	log.Info("请求方法: ", req.Method, "\n请求URL:", req.URL, "\n请求头:", req.Header, "\n请求体:", string(reqbodyBytes),
		"\n响应状态码:", resp.StatusCode, "\n响应头:", resp.Header, "\n响应体:", string(bodyBytes))
	return req, resp, nil
}
