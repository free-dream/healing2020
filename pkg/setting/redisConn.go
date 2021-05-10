package setting

import (
	"github.com/go-redis/redis"
	"healing2020/pkg/tools"
    "strconv"
)

var RedisClient *redis.Client

func init() {
	addr := tools.GetConfig("redis", "addr")
    dbStr := tools.GetConfig("redis", "db")
    var db int
    if dbStr == "" {
        db = 0
    } else {
        db,_ = strconv.Atoi(dbStr)
    }
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     "",
		DB:           db,
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
}

func RedisConn() *redis.Client {
	return RedisClient
}
