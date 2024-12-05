package net

import (
	"extension/web_server"
	"fmt"
	"testing"
)

func TestWeb(t *testing.T) {
	handlers := make(map[string]func(req map[string]any) (resp []byte))
	f := func(req map[string]any) (resp []byte) {
		fmt.Println(req)
		return []byte("Hello World")
	}
	handlers["/test1"] = f
	handlers["/test2"] = f
	web_server.Web(8080, handlers)
}
