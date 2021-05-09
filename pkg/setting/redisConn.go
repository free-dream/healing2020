package setting

import (
	"github.com/go-redis/redis"
	"healing2020/pkg/tools"
)

var RedisClient *redis.Client

func init() {
	addr := tools.GetConfig("redis", "addr")
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     "",
		DB:           0,
		PoolSize:     30,
		MinIdleConns: 10,
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func RedisConnTest() {
	client := RedisConn()
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	client.Set("healing2020:rankCount", 0, 0)
	client.Close()
}

func RedisConn() *redis.Client {
	return RedisClient
}
