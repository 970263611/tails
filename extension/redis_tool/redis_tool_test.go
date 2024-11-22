package redis_tool

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	//args := []string{"-h", "127.0.0.1", "-p", "6379", "-w", "", "-d", "0", "-k", "test"}
	params := make(map[string]any)
	params["addr"] = "127.0.0.1:6379,192.168.227.30:6379"
	params["password"] = ""
	params["db"] = 0

	//key数量
	//params["db_size"] = ""

	//string
	//params["key"] = "test"

	//list
	//params["key"] = "list"
	params["list_index"] = 2

	//zset
	//params["key"] = "zset"
	//params["zset_score"] = "Chinese"

	//hash
	//params["key"] = "hash"
	//params["hash_get"] = "name"

	//set
	//params["key"] = "set"
	//params["set_is"] = " "

	ins := GetInstance(nil)
	bytes := ins.Do(params)
	fmt.Println(bytes)
	fmt.Println(string(bytes))
}
