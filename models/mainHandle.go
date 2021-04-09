package models

import (
    "healing2020/pkg/setting"
    "healing2020/models/statements"
    "encoding/json"
    "database/sql"
    //"strconv"
    //"fmt"

    //"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

type MainMsg struct {
    Sing []SongMsg `json:"sing"`
    Listen []SongMsg `json"listen"`
}

type SongMsg struct {
    User string `json:"user"`
    SendUser string `json:"senduser"`
    Name string `json:"name"`
    Avatar string `json:"avatar"`
    Time string `json:"time"`
    Sex int `json:"sex"`
    More string `json:"more"`

    Id uint `json:"id"`
    Like int `json:"like"`
    Style string `json:"style"`
    Source string `json:"source"`
}

func SendMainMsg() {

}

func LoadSongMsg(sort string,key string) []SongMsg{
    db := setting.MysqlConn()
    defer db.Close()
    var songList []SongMsg = make([]SongMsg,10)
    i := 0

    var rows *sql.Rows
    if key == ""{
        rows,_ = db.Raw("select id,user_id,vod_send,name,praise,source,style,created_at from song order by rand() limit 10").Rows()
    }else {
        if sort == "1" {
            rows,_ = db.Raw("select id,user_id,vod_send,name,praise,source,style,created_at from song where style=? or language=? order by rand() limit 10",key,key).Rows()
        }else {
            rows,_ = db.Raw("select id,user_id,vod_send,name,praise,source,style,created_at from song where style=? or language=? order by created_at,praise desc limit 10",key,key).Rows()
        }
    } 
    defer rows.Close()

    for rows.Next() {
        var song statements.Song
        db.ScanRows(rows,&song)
        songList[i].Id = song.ID
        songList[i].Like = song.Praise
        songList[i].Source = song.Source
        songList[i].Style = song.Style
        userid := song.UserId
        sendid := song.VodSend

        var user statements.User
        db.Model(&statements.User{}).Select("nick_name,sex,avatar").Where("id=?",userid).Find(&user)
        songList[i].User = user.NickName
        songList[i].Sex = user.Sex
        songList[i].Avatar = user.Avatar

        var vod statements.Vod
        db.Model(&statements.Vod{}).Select("name,more").Where("id=?",sendid).Find(&vod)
        songList[i].SendUser = vod.Name
        songList[i].More = vod.More

        i++
    }

    return songList
}

func GetMainMsg(sort string,key string) (MainMsg,error){
    var result MainMsg
    client := setting.RedisConn()
    defer client.Close()
    data1,err1 := client.Get("healing:main:"+key+"SingMsg"+sort).Bytes()
    if err1!=nil {return MainMsg{},err1}
    data2,err2 := client.Get("healing:main:"+key+"ListenMsg"+sort).Bytes()
    if err2!=nil {return MainMsg{},err2}

    var sing []SongMsg
    var listen []SongMsg
    json.Unmarshal(data1,sing)
    json.Unmarshal(data2,listen)
    result.Sing = sing
    result.Listen = listen

    return result,nil
}
