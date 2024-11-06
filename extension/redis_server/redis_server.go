package redis_server

import (
	"basic"
	othertool "basic/tool/other"
	redistool "basic/tool/redis"
	"errors"
	"math"
	"strconv"
)

type RedisServer struct{}

func GetInstance() *RedisServer {
	return &RedisServer{}
}

func (r *RedisServer) GetOrder() int {
	return math.MaxInt64
}

func (r *RedisServer) Register(globalContext *basic.Context) *basic.ComponentMeta {
	p1 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-h",
		StandardName: "host",
		Required:     true,
		CheckMethod: func(s string) error {
			if !othertool.CheckIp(s) {
				errors.New("port is not valid")
			}
			return nil
		},
		Describe: "redis ip地址",
	}
	p2 := basic.Parameter{
		ParamType:    basic.INT,
		CommandName:  "-p",
		StandardName: "port",
		Required:     true,
		CheckMethod: func(s string) error {
			if !othertool.CheckPortByString(s) {
				errors.New("port is not valid")
			}
			return nil
		},
		Describe: "redis端口",
	}
	p3 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-w",
		StandardName: "password",
		Required:     false,
		CheckMethod: func(s string) error {
			if !othertool.CheckPortByString(s) {
				errors.New("port is not valid")
			}
			return nil
		},
		Describe: "redis密码",
	}
	p4 := basic.Parameter{
		ParamType:    basic.INT,
		CommandName:  "-d",
		StandardName: "db",
		Required:     false,
		CheckMethod: func(s string) error {
			if !othertool.CheckPortByString(s) {
				errors.New("db is not valid")
			}
			return nil
		},
		Describe: "redis库号",
	}
	p5 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-k",
		StandardName: "key",
		Required:     true,
		Describe:     "要查询的key值",
	}
	p6 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-zs",
		StandardName: "zset_score",
		Required:     true,
		Describe:     "ZScore():获取元素的score",
	}
	//p7 := basic.Parameter{
	//	ParamType:    basic.INT,
	//	CommandName:  "-li",
	//	StandardName: "list_index",
	//	Required:     true,
	//	Describe:     "LIndex():获取链表下标对应的元素",
	//}
	p7 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-hg",
		StandardName: "hash_get",
		Required:     true,
		Describe:     "HGet():获取某个元素",
	}
	return &basic.ComponentMeta{
		Key:       "redis_server",
		Describe:  "redis连接服务",
		Params:    []basic.Parameter{p1, p2, p3, p4, p5, p6, p7},
		Component: r,
	}
}

func (r *RedisServer) Do(params map[string]any) (resp []byte) {
	hostStr := params["host"].(string)
	poststr := strconv.Itoa(params["port"].(int))
	addrstr := hostStr + ":" + poststr
	passwordStr := params["password"].(string)
	dbStr := params["db"].(int)
	//创建redis连接
	client := redistool.CreateRedisClient(addrstr, passwordStr, dbStr)
	keyStr := params["key"].(string)
	//获取值类型
	typeStr := client.Type(client.Context, keyStr)
	switch typeStr.Val() {
	case "string":
		vel, err := client.Get(client.Context, keyStr).Result()
		if err != nil {
			return []byte("查询异常")
		}
		return []byte(vel)
	case "list":
		list_indexStr := params["list_index"].(int)
		vel, err := client.LIndex(client.Context, keyStr, int64(list_indexStr)).Result()
		if err != nil {
			return []byte("查询异常")
		}
		return []byte(vel)
	case "set":
		//vel, err := client.SMembers(client.Context, keyStr).Result()
		//if err != nil {
		//	return []byte("值为空")
		//}
		//return vel
	case "zset":
		zset_scoreStr := params["zset_score"].(string)
		vel, err := client.ZScore(client.Context, keyStr, zset_scoreStr).Result()
		if err != nil {
			return []byte("查询异常")
		}
		return []byte(strconv.FormatFloat(vel, 'f', -1, 64))
	case "hash":
		hash_getStr := params["hash_get"].(string)
		vel, err := client.HGet(client.Context, keyStr, hash_getStr).Result()
		if err != nil {
			return []byte("查询异常")
		}
		return []byte(vel)
	default:
		return []byte("当前类型暂不支持")
	}
	return []byte("没有匹配的类型")
}
