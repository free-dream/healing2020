package tools

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/garyburd/redigo/redis"
)

type RedisUser struct {
    ID 	uint
}	

func GetUser() RedisUser{
	//连接redis
	r, err := redis.Dial("tcp", GetConfig("redis", "port"))
    if err != nil {
        fmt.Println("Connect to redis error", err)
        os.Exit(1)
    }
    defer r.Close()
    //redis获取数据并绑定json
    var userInf RedisUser
    value, err := redis.Bytes(r.Do("GET", "user"))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    errShal := json.Unmarshal(value, &userInf)
    if errShal != nil {
        fmt.Println(errShal)
    }
    return userInf
}