package redis_server

import (
	"basic"
	redistool "basic/tool/redis"
	"basic/tool/utils"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type RedisServer struct{}

func GetInstance() *RedisServer {
	return &RedisServer{}
}

func (c *RedisServer) GetName() string {
	return "redis_server"
}

func (c *RedisServer) GetDescribe() string {
	return "redis连接服务，例：redis_server -h 127.0.0.1 -p 6379 -k keyStr"
}

func (r *RedisServer) Register(globalContext *basic.Context) *basic.ComponentMeta {
	p1 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-h",
		ConfigName:   "redis.server.ip",
		StandardName: "host",
		Required:     true,
		CheckMethod: func(s string) error {
			if !utils.CheckIp(s) {
				return errors.New("IP不合法")
			}
			return nil
		},
		Describe: "redis ip地址",
	}
	p2 := basic.Parameter{
		ParamType:    basic.INT,
		CommandName:  "-p",
		StandardName: "port",
		ConfigName:   "redis.server.port",
		Required:     true,
		CheckMethod: func(s string) error {
			if !utils.CheckPortByString(s) {
				return errors.New("端口不合法")
			}
			return nil
		},
		Describe: "redis端口",
	}
	p3 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-w",
		StandardName: "password",
		ConfigName:   "redis.server.password",
		Required:     false,
		Describe:     "redis密码",
	}
	p4 := basic.Parameter{
		ParamType:    basic.INT,
		CommandName:  "-d",
		StandardName: "db",
		ConfigName:   "redis.server.db",
		Required:     false,
		Describe:     "redis库号",
	}
	p5 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-k",
		StandardName: "key",
		Required:     false,
		Describe:     "要查询的key值，例：redis_server -h 127.0.0.1 -p 6379 -k keyStr",
	}
	p6 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-Z",
		StandardName: "zset_score",
		Required:     false,
		Describe:     "获取某元素的score，例：redis_server -h 127.0.0.1 -p 6379 -k zset -Z zsetScore",
	}
	p7 := basic.Parameter{
		ParamType:    basic.INT,
		CommandName:  "-L",
		StandardName: "list_index",
		Required:     false,
		Describe:     "获取链表下标对应的元素，例：redis_server -h 127.0.0.1 -p 6379 -k list -L list_index",
	}
	p8 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-H",
		StandardName: "hash_get",
		Required:     false,
		Describe:     "获取某个元素，例：redis_server -h 127.0.0.1 -p 6379 -k hash -H name",
	}
	p9 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-S",
		StandardName: "set_is",
		Required:     false,
		Describe:     "判断元素是否在集合中，例：redis_server -h 127.0.0.1 -p 6379 -k set -S abc",
	}
	p10 := basic.Parameter{
		ParamType:    basic.NO_VALUE,
		CommandName:  "-D",
		StandardName: "db_size",
		Required:     false,
		Describe:     "查看当前数据库key的数量，例：redis_server -h 127.0.0.1 -p 6379 -D",
	}

	return &basic.ComponentMeta{
		ComponentType: basic.EXECUTE,
		Params:        []basic.Parameter{p1, p2, p3, p4, p5, p6, p7, p8, p9, p10},
		Component:     r,
	}
}
func (r *RedisServer) Start(globalContext *basic.Context) error {
	return nil
}

func (r *RedisServer) Do(params map[string]any) (resp []byte) {
	//创建redis连接
	hostStr := params["host"].(string)
	poststr := strconv.Itoa(params["port"].(int))
	addrstr := hostStr + ":" + poststr
	passwordStr, flag := params["password"].(string)
	if !flag {
		passwordStr = ""
	}
	dbStr, true := params["db"].(int)
	if !true {
		dbStr = 0
	}
	var client *redistool.RedisClient
	client, err := redistool.CreateRedisClient(addrstr, passwordStr, dbStr)
	if err != nil {
		return []byte("redis连接失败")
	}

	//查看当前数据库key的数量
	_, ok := params["db_size"]
	if ok {
		vel, err := client.DBSize(client.Context).Result()
		if err != nil {
			log.Error("-D参数解析结果为空", err.Error())
			return []byte("-D参数解析结果为空" + err.Error())
		}
		return []byte(strconv.FormatInt(vel, 10))
	}

	//获取key值
	keyStr, _ := params["key"].(string)
	if len(keyStr) == 0 {
		log.Error("-K参数解析结果为空,key值为：", keyStr)
		return []byte("-k参数解析结果为空,key值为：" + keyStr)
	}

	//获取值类型
	typeStr := client.Type(client.Context, keyStr)
	switch typeStr.Val() {
	case "string":
		vel, err := client.Get(client.Context, keyStr).Result()
		if err != nil {
			log.Error("查询异常：", keyStr)
			return []byte("查询异常：" + keyStr)
		}
		return []byte(vel)
	case "list":
		//1.获取链表下标对应的元素
		list_indexStr, _ := params["list_index"].(int)
		if list_indexStr > 0 {
			vel, err := client.LIndex(client.Context, keyStr, int64(list_indexStr)).Result()
			if err != nil {
				log.Error("-L参数解析结果为空,key值为：", keyStr)
				return []byte("-L参数解析结果为空,key值为：" + keyStr)
			}
			return []byte(vel)
		}

		//默认获取所有元素
		length, _ := client.LLen(client.Context, keyStr).Result()
		if length > 0 {
			vel, err := client.LRange(client.Context, keyStr, 0, length-1).Result()
			if err != nil {
				log.Error("查询异常：", keyStr)
				return []byte("查询异常：" + keyStr)
			}
			data, err := json.Marshal(vel)
			if err != nil {
				log.Error("类型转换失败：", keyStr)
				return []byte("类型转换失败：" + keyStr)
			}
			return data
		}
	case "set":
		//1.判断元素是否在集合中
		set_isStr, _ := params["set_is"].(string)
		if len(set_isStr) > 0 {
			vel, err := client.SIsMember(client.Context, keyStr, set_isStr).Result()
			if err != nil {
				log.Error("-S参数解析结果为空,key值为：", keyStr)
				return []byte("-S参数解析结果为空,key值为：" + keyStr)
			}
			if vel {
				return []byte("true")
			} else {
				return []byte("false")
			}
		}

		//默认获取所有元素
		vel, err := client.SMembers(client.Context, keyStr).Result()
		if err != nil {
			log.Error("查询异常：", keyStr)
			return []byte("查询异常：" + keyStr)
		}
		data, err := json.Marshal(vel)
		if err != nil {
			log.Error("类型转换失败：", keyStr)
			return []byte("类型转换失败：" + keyStr)
		}
		return data
	case "zset":
		//1.获取某元素的score
		zset_scoreStr, _ := params["zset_score"].(string)
		if len(zset_scoreStr) > 0 {
			vel, err := client.ZScore(client.Context, keyStr, zset_scoreStr).Result()
			if err != nil {
				log.Error("-Z参数解析结果为空,key值为：", keyStr)
				return []byte("-Z参数解析结果为空,key值为：" + keyStr)
			}
			return []byte(strconv.FormatFloat(vel, 'f', -1, 64))
		}

		//默认获取所有元素
		count, _ := client.ZCard(client.Context, keyStr).Result()
		if count > 0 {
			vel, err := client.ZRange(client.Context, keyStr, 0, count-1).Result()
			if err != nil {
				log.Error("查询异常：", keyStr)
				return []byte("查询异常：" + keyStr)
			}
			data, err := json.Marshal(vel)
			if err != nil {
				log.Error("类型转换失败：", keyStr)
				return []byte("类型转换失败：" + keyStr)
			}
			return data
		}

	case "hash":
		//1.获取某个元素
		hash_getStr, _ := params["hash_get"].(string)
		if len(hash_getStr) > 0 {
			vel, err := client.HGet(client.Context, keyStr, hash_getStr).Result()
			if err != nil {
				log.Error("-H参数解析结果为空,key值为：", keyStr)
				return []byte("-H参数解析结果为空,key值为：" + keyStr)
			}
			return []byte(vel)
		}

		//默认获取全部元素
		vel, err := client.HGetAll(client.Context, keyStr).Result()
		if err != nil {
			log.Error("查询异常：", keyStr)
			return []byte("查询异常：" + keyStr)
		}
		data, err := json.Marshal(vel)
		if err != nil {
			log.Error("类型转换失败：", keyStr)
			return []byte("类型转换失败：" + keyStr)
		}
		return data
	default:
		log.Error("当前key不存在或类型暂不支持解析", keyStr)
		return []byte("当前key不存在或类型暂不支持解析" + keyStr)
	}
	log.Error("当前key不存在或类型暂不支持解析", keyStr)
	return []byte("当前key不存在或类型暂不支持解析" + keyStr)
}
