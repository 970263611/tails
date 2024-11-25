package redistool

import (
	"context"
	redis "github.com/redis/go-redis/v9"
	"strings"
)

type RedisClient struct {
	redis.Cmdable
	context.Context
}

/*
*
地址 127.0.0.1:6379 ; 密码 没有送"" ; 库号
创建redis_client=单点模式
*/
func CreateRedisClientSingle(addr, password string, db int) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	//return &RedisClient{rdb, nil, nil, context.Background()}, nil
	return &RedisClient{rdb, context.Background()}, nil
}

/*
*
地址 127.0.0.1:6379,127.0.0.1:6378 ; 密码 没有送"" ;
创建redis_client=集群模式
*/
func CreateRedisClientCluster(addr string, password string) (*RedisClient, error) {
	parts := strings.Split(addr, ",")
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    parts,
		Password: password,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	//return &RedisClient{nil, nil, rdb, context.Background()}, nil
	return &RedisClient{rdb, context.Background()}, nil
}

/*
*
地址 127.0.0.1:6379,127.0.0.1:6378 ; 密码 没有送"" ;
创建redis_client=哨兵模式
*/
func CreateRedisClientSentinel(sentinelAddrs string, password string, db int, sentinelPassword string, masterName string) (*RedisClient, error) {
	parts := strings.Split(sentinelAddrs, ",")
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs:    parts,            //哨兵地址
		SentinelPassword: sentinelPassword, //哨兵密码
		MasterName:       masterName,       //主节点的名称
		Password:         password,         //节点密码
		DB:               db,               //库号
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return &RedisClient{rdb, context.Background()}, nil
}
