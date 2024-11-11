package net

import (
	"fmt"
	"log"
	"net/url"
	"testing"
)

// ResponseData 是一个示例结构体，用于解码 HTTP 响应体
type ResponseData struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func TestDoGetToString(t *testing.T) {
	// URL
	urlStr := "https://jsonplaceholder.typicode.com/posts"
	// 查询参数
	queryParams := url.Values{}
	queryParams.Add("userId", "1")

	// 发送 GET 请求
	toString, err := Get(urlStr, queryParams)
	if err != nil {
		log.Fatalf("GET 请求失败: %v", err)
	}
	fmt.Printf("GET 请求响应: %+v\n", toString)
}

func TestGetToStruct(t *testing.T) {
	// URL
	urlStr := "https://jsonplaceholder.typicode.com/posts"
	// 查询参数
	queryParams := url.Values{}
	queryParams.Add("userId", "1")
	// 用于接收响应的结构体实例
	var responseData ResponseData
	// 发送 GET 请求
	err := GetToStruct(urlStr, queryParams, &responseData)
	if err != nil {
		log.Fatalf("GET 请求失败: %v", err)
	}
	fmt.Printf("GET 请求响应: %+v\n", responseData)
}
