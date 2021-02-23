package controller

import (
    "encoding/json"
    "fmt"

    "healing2020/models/statements"
    "healing2020/pkg/tools"

    "github.com/garyburd/redigo/redis"
)

func GetRedisUser() statements.User {
	addr := tools.GetConfig("redis", "addr")
	//连接redis
	r, err := redis.Dial("tcp", addr)
    if err != nil {
        fmt.Println("Connect to redis error", err)    
    }
    defer r.Close()
    //redis获取数据并绑定json
    var userInf statements.User
    value, err := redis.Bytes(r.Do("GET", "user"))
    if err != nil {
        fmt.Println(err)
    }
    errShal := json.Unmarshal(value, &userInf)
    if errShal != nil {
        fmt.Println(errShal)
    }
    return userInf
}