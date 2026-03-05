package dao

import (
	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

// InitRedis
func InitRedis() *redis.Client {
	Rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // Redis 地址
		Password: "",               // 密码
		DB:       0,                // 使用默认数据库
	})
	return Rdb
}
