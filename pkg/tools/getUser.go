package tools

import (
    "encoding/json"
    "fmt"

    "github.com/garyburd/redigo/redis"
    "github.com/jinzhu/gorm"
)

type RedisUser struct {
    gorm.Model
    OpenId string
    NickName string
    TrueName string
    More string
    Campus string
    Avatar string
    Phone string
    Sex int
    Hobby string
    Money int
    Setting1 int
    Setting2 int
    Setting3 int
}

func GetUser() RedisUser {
	addr := GetConfig("redis", "addr")
	//连接redis
	r, err := redis.Dial("tcp", addr)
    if err != nil {
        fmt.Println("Connect to redis error", err)    
    }
    defer r.Close()
    //redis获取数据并绑定json
    var userInf RedisUser
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