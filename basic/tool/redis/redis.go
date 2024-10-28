package redistool

import (
	"context"
	redis "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
	context.Context
}

/*
*
地址 127.0.0.1:6379 ; 密码 没有送"" ; 库号
*/
func CreateRedisClient(addr, password string, db int) *RedisClient {
	option := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	}
	return CreateRedisClientByOptions(option)
}

/*
*
通过Option创建redis_client
*/
func CreateRedisClientByOptions(opt *redis.Options) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisClient{rdb, context.Background()}
}
