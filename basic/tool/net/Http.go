package net

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}
	return string(body), nil
}

func Post(url string, data string) (string, error) {
	param := []byte(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(param))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}
	return string(body), nil
}
