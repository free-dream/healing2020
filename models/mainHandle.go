package models

import (
    "healing2020/pkg/setting"
    "healing2020/models/statements"
    "encoding/json"
    "database/sql"
    //"strconv"
    //"fmt"

    "github.com/jinzhu/gorm"
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
    Singer string `json:"singer"`
}

func SendMainMsg() {
    client := setting.RedisConn()
    defer client.Close()
    
    var sortArr = []string{"0","1"}
    var keyArr = []string{"","ACG","流行","古风","民谣","摇滚","抖音热搜","国语","英语","日语","粤语"}
    for _,sort := range sortArr {
        for _,key :=  range keyArr {
            listen := LoadSongMsg(sort,key)
            sing := LoadVodMsg(sort,key)

            keyname := "healing2020:Main:"+key+"ListenMsg"+sort
            client.Set(keyname,listen,0)
            keyname = "healing2020:Main:"+key+"Sing"+sort
            client.Set(keyname,sing,0)
        }
    }
}

func LoadSongMsg(sort string,key string) []SongMsg{
    db := setting.MysqlConn()
    defer db.Close()
    var songList []SongMsg = make([]SongMsg,10)
    i := 0

    var rows *sql.Rows
    var result *gorm.DB
    if key == ""{
        result = db.Raw("select id,user_id,vod_send,name,praise,source,style,created_at from song order by rand() limit 10")
        if isListNil(result) {
            return nil
        }
        rows,_ = result.Rows()
    }else {
        if sort == "1" {
            result = db.Raw("select id,user_id,vod_send,name,praise,source,style,created_at from song where style=? or language=? order by rand() limit 10",key,key)
            if isListNil(result) {
                return nil
            }
            rows,_ = result.Rows()
        }else {
            result = db.Raw("select id,user_id,vod_send,name,praise,source,style,created_at from song where style=? or language=? order by created_at,praise desc limit 10",key,key)
            if isListNil(result) {
                return nil
            }
            rows,_ = result.Rows()
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

func LoadVodMsg(sort string,key string) []SongMsg{
    db := setting.MysqlConn()
    defer db.Close()
    var vodList []SongMsg = make([]SongMsg,10)
    i := 0

    var rows *sql.Rows
    var result *gorm.DB
    if key == ""{
        result = db.Raw("select id,user_id,name,singer,more,created_at from vod order by rand() limit 10")
        rows,_ = result.Rows()
    }else {
        if sort == "1" {
            result = db.Raw("select id,user_id,name,singer,more,created_at from vod where style=? or language=? order by rand() limit 10",key,key)
            if isListNil(result) {
                return nil
            }
            rows,_ = result.Rows()
        }else {
            result = db.Raw("select id,user_id,name,singer,more,created_at from vod where style=? or language=? order by created_at,praise desc limit 10",key,key)
            if isListNil(result) {
                return nil
            }
            rows,_ = result.Rows()
        }
    } 
    defer rows.Close()

    for rows.Next() {
        var vod statements.Vod
        db.ScanRows(rows,&vod)
        vodList[i].Id = vod.ID
        vodList[i].Name = vod.Name
        vodList[i].More = vod.More
        //vodList[i].Time = vod.CreatedAt
        vodList[i].Singer = vod.Singer
        
        userid := vod.UserId

        var user statements.User
        db.Model(&statements.User{}).Select("nick_name,sex,avatar").Where("id=?",userid).Find(&user)
        vodList[i].User = user.NickName
        vodList[i].Sex = user.Sex
        vodList[i].Avatar = user.Avatar

        i++
    }

    return vodList
}

func GetMainMsg(sort string,key string) (MainMsg,error){
    var result MainMsg
    client := setting.RedisConn()
    defer client.Close()
    data1,err1 := client.Get("healing2020:Main:"+key+"SingMsg"+sort).Bytes()
    if err1!=nil {return MainMsg{},err1}
    data2,err2 := client.Get("healing2020:Main:"+key+"ListenMsg"+sort).Bytes()
    if err2!=nil {return MainMsg{},err2}

    var sing []SongMsg
    var listen []SongMsg
    json.Unmarshal(data1,sing)
    json.Unmarshal(data2,listen)
    result.Sing = sing
    result.Listen = listen

    return result,nil
}

func isListNil(result *gorm.DB) bool{
    count := result.RowsAffected
    if count == 0 {
        return true
    }
    return false
}
