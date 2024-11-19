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
func CreateRedisClient(addr, password string, db int) (*RedisClient, error) {
	option := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	}
	options, err := CreateRedisClientByOptions(option)
	if err != nil {
		return nil, err
	}
	return options, nil
}

/*
*
通过Option创建redis_client
*/
func CreateRedisClientByOptions(opt *redis.Options) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password, // no password set
		DB:       opt.DB,       // use default DB
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return &RedisClient{rdb, context.Background()}, nil
}
