package redis_server

import (
	cons "basic/constants"
	iface "basic/interfaces"
	redistool "basic/tool/redis"
	"basic/tool/utils"
	"encoding/json"
	"errors"
	"strconv"
)

type RedisServer struct{}

func GetInstance(globalContext iface.Context) iface.Component {
	return &RedisServer{}
}

func (c *RedisServer) GetName() string {
	return "redis_server"
}

func (c *RedisServer) GetDescribe() string {
	return "redis连接服务"
}

func (r *RedisServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, "-h", "redis.server.host", "host", true,
		func(s string) error {
			if !utils.CheckIp(s) {
				return errors.New("IP不合法")
			}
			return nil
		}, "redis ip地址")
	cm.AddParameters(cons.INT, "-p", "redis.server.port", "port", true,
		func(s string) error {
			if !utils.CheckPortByString(s) {
				return errors.New("端口不合法")
			}
			return nil
		}, "redis端口")
	cm.AddParameters(cons.STRING, "-w", "", "password", false, nil, "redis密码")
	cm.AddParameters(cons.INT, "-d", "", "db", false, nil, "redis库号")
	cm.AddParameters(cons.STRING, "-k", "", "key", false, nil, "要查询的key值")
	cm.AddParameters(cons.STRING, "-Z", "", "zset_score", false, nil, "ZScore():获取某元素的score")
	cm.AddParameters(cons.INT, "-L", "", "list_index", false, nil, "LIndex():获取链表下标对应的元素")
	cm.AddParameters(cons.STRING, "-H", "", "hash_get", false, nil, "HGet():获取某个元素")
	cm.AddParameters(cons.STRING, "-S", "", "set_is", false, nil, "SIsMember():判断元素是否在集合中")
	cm.AddParameters(cons.NO_VALUE, "-D", "", "db_size", false, nil, "DBSize():查看当前数据库key的数量")
}
func (r *RedisServer) Start() error {
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
	client = redistool.CreateRedisClient(addrstr, passwordStr, dbStr)

	//查看当前数据库key的数量
	_, ok := params["db_size"]
	if ok {
		vel, err := client.DBSize(client.Context).Result()
		if err != nil {
			return []byte("-D参数解析结果为空")
		}
		return []byte(strconv.FormatInt(vel, 10))
	}

	//获取key值
	keyStr, _ := params["key"].(string)
	if len(keyStr) == 0 {
		return []byte("-k参数解析结果为空")
	}

	//获取值类型
	typeStr := client.Type(client.Context, keyStr)
	switch typeStr.Val() {
	case "string":
		vel, err := client.Get(client.Context, keyStr).Result()
		if err != nil {
			return []byte("查询异常：" + keyStr)
		}
		return []byte(vel)
	case "list":
		//1.获取链表下标对应的元素
		list_indexStr, _ := params["list_index"].(int)
		if list_indexStr > 0 {
			vel, err := client.LIndex(client.Context, keyStr, int64(list_indexStr)).Result()
			if err != nil {
				return []byte("-L参数解析结果为空,key值为：" + keyStr)
			}
			return []byte(vel)
		}

		//默认获取所有元素
		length, _ := client.LLen(client.Context, keyStr).Result()
		if length > 0 {
			vel, err := client.LRange(client.Context, keyStr, 0, length-1).Result()
			if err != nil {
				return []byte("查询异常：" + keyStr)
			}
			data, err := json.Marshal(vel)
			if err != nil {
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
			return []byte("查询异常：" + keyStr)
		}
		data, err := json.Marshal(vel)
		if err != nil {
			return []byte("类型转换失败：" + keyStr)
		}
		return data
	case "zset":
		//1.获取某元素的score
		zset_scoreStr, _ := params["zset_score"].(string)
		if len(zset_scoreStr) > 0 {
			vel, err := client.ZScore(client.Context, keyStr, zset_scoreStr).Result()
			if err != nil {
				return []byte("-Z参数解析结果为空,key值为：" + keyStr)
			}
			return []byte(strconv.FormatFloat(vel, 'f', -1, 64))
		}

		//默认获取所有元素
		count, _ := client.ZCard(client.Context, keyStr).Result()
		if count > 0 {
			vel, err := client.ZRange(client.Context, keyStr, 0, count-1).Result()
			if err != nil {
				return []byte("查询异常：" + keyStr)
			}
			data, err := json.Marshal(vel)
			if err != nil {
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
				return []byte("-H参数解析结果为空,key值为：" + keyStr)
			}
			return []byte(vel)
		}

		//默认获取全部元素
		vel, err := client.HGetAll(client.Context, keyStr).Result()
		if err != nil {
			return []byte("查询异常：" + keyStr)
		}
		data, err := json.Marshal(vel)
		if err != nil {
			return []byte("类型转换失败：" + keyStr)
		}
		return data
	default:
		return []byte("当前类型暂不支持解析")
	}
	return []byte("当前key不存")
}
