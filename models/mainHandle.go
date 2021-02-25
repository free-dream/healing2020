package models

import (
    "healing2020/pkg/setting"
    "encoding/json"
    //"strconv"
)

type MainMsg struct {
    Sing []SongMsg `json:"sing"`
    Listen []SongMsg `json"listen"`
}

type SongMsg struct {
    User string `json:"user"`
    Name string `json:"name"`
    Avatar string `json:"avatar"`
    Time string `json:"time"`
    Sex int `json:"sex"`
    More string `json:"more"`

    Id uint `json:"id"`
    Like int `json:"like"`
    Style string `json:"style"`
}

func SendMainMsg() {

}

func GetMainMsg(sort string,key string) (MainMsg,error){
    var result MainMsg
    client := setting.RedisConn()
    defer client.Close()
    data1,err1 := client.Get(key+"SingMsg"+sort).Bytes()
    if err1!=nil {return MainMsg{},err1}
    data2,err2 := client.Get(key+"ListenMsg"+sort).Bytes()
    if err2!=nil {return MainMsg{},err2}

    var sing []SongMsg
    var listen []SongMsg
    json.Unmarshal(data1,sing)
    json.Unmarshal(data2,listen)
    result.Sing = sing
    result.Listen = listen

    return result,nil
}
