package redis_tool

import (
	cons "basic/constants"
	iface "basic/interfaces"
	redistool "basic/tool/redis"
	"basic/tool/utils"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type RedisServer struct{}

func GetInstance(globalContext iface.Context) iface.Component {
	return &RedisServer{}
}

func (c *RedisServer) GetName() string {
	return "redis_tool"
}

func (c *RedisServer) GetDescribe() string {
	return "redis连接服务，例：redis_tool -addr 127.0.0.1:6379 -k keyStr"
}

func (r *RedisServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, cons.LOWER_A, "addr", "addr", true,
		func(s string) error {
			if !utils.CheckAddr(s) {
				return errors.New("主机地址不合法")
			}
			return nil
		}, "redis主机地址")
	cm.AddParameters(cons.STRING, cons.LOWER_W, "password", "password", false, nil, "redis密码")
	cm.AddParameters(cons.INT, cons.LOWER_N, "db", "db", false, nil, "redis库号")
	cm.AddParameters(cons.STRING, cons.LOWER_K, "", "key", false, nil, "要查询的key值，例：redis_server -addr 127.0.0.1:6379 -k keyStr")
	cm.AddParameters(cons.STRING, cons.UPPER_Z, "", "zset_score", false, nil, "获取某元素的score，例：redis_server -addr 127.0.0.1:6379 -k zset -Z zsetScore")
	cm.AddParameters(cons.INT, cons.UPPER_L, "", "list_index", false, nil, "获取链表下标对应的元素，例：redis_server -addr 127.0.0.1:6379 -k list -L list_index")
	cm.AddParameters(cons.STRING, cons.UPPER_H, "", "hash_get", false, nil, "获取某个元素，例：redis_server -addr 127.0.0.1:6379 -k hash -H name")
	cm.AddParameters(cons.STRING, cons.UPPER_S, "", "set_is", false, nil, "判断元素是否在集合中，例：redis_server -addr 127.0.0.1:6379 -k set -S abc")
	cm.AddParameters(cons.NO_VALUE, cons.UPPER_D, "", "db_size", false, nil, "查看当前数据库key的数量，例：redis_server -addr 127.0.0.1:6379 -D")
}

func (r *RedisServer) Do(params map[string]any) (resp []byte) {
	//创建redis连接
	addrstr := params["addr"].(string)
	passwordStr, flag := params["password"].(string)
	if !flag {
		passwordStr = ""
	}
	dbStr, true := params["db"].(int)
	if !true {
		dbStr = 0
	}
	var client *redistool.RedisClient
	errorList := make([]string, 0) // 用于收集错误信息的切片
	parts := strings.Split(addrstr, ",")
	for _, part := range parts {
		fmt.Println(part)
		c, err := redistool.CreateRedisClient(part, passwordStr, dbStr)
		if err != nil {
			errorList = append(errorList, fmt.Sprintf("创建Redis客户端失败，地址: %s，错误: %v", part, err))
			continue
		}
		client = c
		// 如果成功创建了客户端，就跳出循环，不需要继续尝试其他地址了
		break
	}
	if client == nil {
		return []byte(strings.Join(errorList, "\n"))
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
