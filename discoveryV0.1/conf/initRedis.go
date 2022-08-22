package conf

import "github.com/go-redis/redis"

// redis

// 定义一个全局变量
var redisdb *redis.Client

func initRedis() (err error) {
	redisdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.75.102:6379", // 指定
		Password: "",
		DB:       0, // redis一共16个库，指定其中一个库即可
	})
	_, err = redisdb.Ping().Result()
	return
}
