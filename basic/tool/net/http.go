package net

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
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
	return resp, nil
}
