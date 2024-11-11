package net

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Get(urlStr string, params url.Values) (string, error) {
	// 构建完整的 URL（如果有查询参数）
	var fullURL string
	if params != nil {
		u, err := url.Parse(urlStr)
		if err != nil {
			return "", fmt.Errorf("解析 URL 失败: %v", err)
		}
		u.RawQuery = params.Encode()
		fullURL = u.String()
	} else {
		fullURL = urlStr
	}
	//发送请求
	resp, err := http.Get(fullURL)
	if err != nil {
		log.Error("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response body:", err)
		return "", err
	}
	return string(body), nil
}

func GetToStruct(urlStr string, params url.Values, response interface{}) error {
	// 构建完整的 URL（如果有查询参数）
	var fullURL string
	if params != nil {
		u, err := url.Parse(urlStr)
		if err != nil {
			return fmt.Errorf("解析 URL 失败: %v", err)
		}
		u.RawQuery = params.Encode()
		fullURL = u.String()
	} else {
		fullURL = urlStr
	}
	//发送请求
	resp, err := http.Get(fullURL)
	if err != nil {
		log.Error("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()
	// 读取响应体
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %v", err)
	}
	// 解码响应体到结构体
	err = json.Unmarshal(bodyBytes, response)
	if err != nil {
		return fmt.Errorf("解码响应体失败: %v", err)
	}
	return nil
}

func Post(urlStr string, postData interface{}) (string, error) {
	// 创建请求体
	var reqBody *strings.Reader
	if postData != nil {
		postDataBytes, err := json.Marshal(postData)
		if err != nil {
			return "", fmt.Errorf("编码 POST 数据失败: %v", err)
		}
		reqBody = strings.NewReader(string(postDataBytes))
	} else {
		return "", fmt.Errorf("post请求请求体不能为空")
	}
	// 创建请求
	req, err := http.NewRequest(http.MethodPost, urlStr, reqBody)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()
	// 读取响应体
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %v", err)
	}
	return string(bodyBytes), nil
}

func PostToStruct(urlStr string, postData interface{}, response interface{}) error {
	// 创建请求体
	var reqBody *strings.Reader
	if postData != nil {
		postDataBytes, err := json.Marshal(postData)
		if err != nil {
			return fmt.Errorf("编码 POST 数据失败: %v", err)
		}
		reqBody = strings.NewReader(string(postDataBytes))
	}
	// 创建请求
	req, err := http.NewRequest(http.MethodPost, urlStr, reqBody)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()
	// 读取响应体
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %v", err)
	}
	// 解码响应体到结构体
	err = json.Unmarshal(bodyBytes, response)
	if err != nil {
		return fmt.Errorf("解码响应体失败: %v", err)
	}
	return nil
}
