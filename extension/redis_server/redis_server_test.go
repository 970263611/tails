package redis_server

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	//args := []string{"-h", "127.0.0.1", "-p", "6379", "-w", "", "-d", "0", "-k", "test"}
	params := make(map[string]any)
	params["host"] = "127.0.0.1"
	params["port"] = 6379
	params["password"] = ""
	params["db"] = 0
	//string
	//params["key"] = "test"

	//list
	//params["key"] = "list"
	//params["list_index"] = 2

	//zset
	//params["key"] = "zset"
	//params["zset_score"] = "Chinese"

	//hash
	params["key"] = "hash"
	params["hash_get"] = "name"

	ins := GetInstance()
	bytes := ins.Do(params)
	fmt.Println(bytes)
	fmt.Println(string(bytes))
}
