package setting

import (
    "healing2020/pkg/tools"
    "github.com/go-redis/redis"
)

func RedisConnTest() {
    client := RedisConn()
    _,err := client.Ping().Result()
    if err != nil {
        panic(err)
    }
    client.Set("healing:rankCount",0,0)
    client.Close()
}

func RedisConn() *redis.Client{
    addr := tools.GetConfig("redis","addr")
    client := redis.NewClient(&redis.Options{
        Addr: addr,
        Password: "",
        DB: 0,
    })
    return client
}
