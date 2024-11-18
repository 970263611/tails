package net

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"
)

// 发送数据的结构体
type PostData struct {
	username string
	password string
}

// 接收数据的结构体
type ResponseData struct {
	Code    int
	Data    []DataEntry
	Message string
}
type DataEntry struct {
	CenterFlg string
	TotalNum  int
}

type DataEntry02 struct {
	Ddd      string
	TotalNum int
}

func TestGetRespString(t *testing.T) {
	// URL
	urlStr := "http://localhost:9999/zhhGetTest"
	// 参数
	queryParams := url.Values{}
	queryParams.Add("userId", "1")
	// Header 如没有必要值header可为null
	header := http.Header{}
	header.Set("Authorization", "dsadasd13213214dssafsadfdsf")
	// 发送 GET 请求
	resp, err := GetRespString(urlStr, nil, header)
	if err != nil {
		log.Fatalf("GET 请求失败: %v", err)
	}
	fmt.Printf("GET 请求响应: %+v\n", resp)
}

func TestGetRespStruct(t *testing.T) {
	// URL
	urlStr := "http://localhost:9999/zhhGetTest"
	// 参数
	queryParams := url.Values{}
	queryParams.Add("userId", "1")
	// Header 如没有必要值header可为null
	header := http.Header{}
	header.Set("Authorization", "dsadasd13213214dssafsadfdsf")
	// 用于接收响应的结构体实例
	var responseData ResponseData
	// 发送 GET 请求
	err := GetRespStruct(urlStr, nil, header, &responseData)
	if err != nil {
		log.Fatalf("GET 请求失败: %v", err)
	}
	entries := responseData.Data
	for _, entry := range entries {
		// 打印每个结构体的字段
		fmt.Printf("CenterFlg: %s, TotalNum: %d\n", entry.CenterFlg, entry.TotalNum)
	}

}

func TestPostRespString(t *testing.T) {
	// URL
	urlStr := "http://localhost:9999/zhhTest"
	data := PostData{
		username: "value1",
		password: "1234455",
	}
	// Header 如没有必要值header可为null
	header := http.Header{}
	header.Set("Authorization", "dsadasd13213214dssafsadfdsf")
	// 发送 POST 请求
	resp, err := PostRespString(urlStr, data, header)
	if err != nil {
		log.Fatalf("POST 请求失败: %v", err)
	}
	fmt.Printf("POST 请求响应: %+v\n", resp)
}

func TestPostRespStruct(t *testing.T) {
	// URL
	urlStr := "http://localhost:9999/zhhTest"
	data := PostData{
		username: "value1",
		password: "13243",
	}
	// Header 如没有必要值header可为null
	header := http.Header{}
	header.Set("Authorization", "dsadasd13213214dssafsadfdsf")
	// 用于接收响应的结构体实例
	var responseData ResponseData
	// 发送 POST 请求
	err := PostRespStruct(urlStr, data, header, &responseData)
	if err != nil {
		log.Fatalf("POST 请求失败: %v", err)
	}
	entries := responseData.Data
	for _, entry := range entries {
		// 打印每个结构体的字段
		fmt.Printf("CenterFlg: %s, TotalNum: %d\n", entry.CenterFlg, entry.TotalNum)
	}
}

func TestPutRespString(t *testing.T) {
	// URL
	urlStr := "http://localhost:9999/zhhPutTest"
	// 参数
	queryParams := url.Values{}
	queryParams.Add("userId", "1")
	// Header 如没有必要值header可为null
	header := http.Header{}
	header.Set("Authorization", "dsadasd13213214dssafsadfdsf")
	// 发送 GET 请求
	resp, err := PutRespString(urlStr, nil, nil, header)
	if err != nil {
		log.Fatalf("PUT 请求失败: %v", err)
	}
	fmt.Printf("PUT 请求响应: %+v\n", resp)
}

func TestPutRespStruct(t *testing.T) {
	// URL
	urlStr := "http://localhost:9999/zhhPutTest"
	// 参数
	queryParams := url.Values{}
	queryParams.Add("userId", "1")
	// Header 如没有必要值header可为null
	header := http.Header{}
	header.Set("Authorization", "dsadasd13213214dssafsadfdsf")
	// 用于接收响应的结构体实例
	var responseData ResponseData
	// 发送 GET 请求
	err := PutRespStruct(urlStr, nil, nil, header, &responseData)
	if err != nil {
		log.Fatalf("PUT 请求失败: %v", err)
	}
	entries := responseData.Data
	for _, entry := range entries {
		// 打印每个结构体的字段
		fmt.Printf("CenterFlg: %s, TotalNum: %d\n", entry.CenterFlg, entry.TotalNum)
	}

}
