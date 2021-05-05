package models

import (
    "healing2020/pkg/setting"
    "healing2020/models/statements"
    "encoding/json"
    "database/sql"
    //"strconv"
    //"fmt"
    "errors"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"

    "time"
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
    Time time.Time `json:"time"`
    Sex int `json:"sex"`
    More string `json:"more"`

    Id uint `json:"id"`
    Like int `json:"like"`
    Style string `json:"style"`
    Source string `json:"source"`
    Singer string `json:"singer"`
    UserId uint `json:"userid"`
}

func SendMainMsg() {
    client := setting.RedisConn()
    defer client.Close()
    
    var sortArr = []string{"0","1"}
    var keyArr = []string{"","ACG","流行","古风","民谣","摇滚","抖音热歌","其他","国语","英语","日语","粤语"}
    for _,sort := range sortArr {
        for _,key :=  range keyArr {
            listenRaw := LoadSongMsg(sort,key)
            singRaw := LoadVodMsg(sort,key)
            listen,_ := json.Marshal(listenRaw)
            sing,_ := json.Marshal(singRaw)

            keyname := "healing2020:Main:"+key+"ListenMsg"+sort
            client.Set(keyname,listen,0)
            keyname = "healing2020:Main:"+key+"SingMsg"+sort
            client.Set(keyname,sing,0)
            //fmt.Println(result.Err())
        }
    }
}

func LoadSongMsg(sort string,key string) []SongMsg{
    db := setting.MysqlConn()
    defer db.Close()
    var songList []SongMsg = make([]SongMsg,8)
    i := 0

    var rows *sql.Rows
    var result *gorm.DB
    if key == ""{
        result = db.Raw("select id,user_id,vod_send,name,praise,source,style,created_at from song order by rand() limit 30")
        rows,_ = result.Rows()
    }else {
        if sort == "1" {
            result = db.Raw("select id,user_id,vod_send,name,praise,source,style,created_at from song where style=? or language=? order by rand() limit 30",key,key)
            rows,_ = result.Rows()
        }else {
            result = db.Raw("select id,user_id,vod_send,name,praise,source,style,created_at from song where style=? or language=? order by created_at,praise desc limit 30",key,key)
            rows,_ = result.Rows()
        }
    } 

    if rows == nil {
        return songList
    }

    defer rows.Close()

    for rows.Next() {
        var song statements.Song
        db.ScanRows(rows,&song)
        songList[i].Id = song.ID
        songList[i].Like = song.Praise
        songList[i].Source = song.Source
        songList[i].Style = song.Style
        songList[i].Time = song.CreatedAt
        userid := song.UserId
        sendid := song.VodSend

        var user statements.User
        db.Model(&statements.User{}).Select("nick_name,sex,avatar").Where("id=?",userid).Find(&user)
        songList[i].User = user.NickName
        songList[i].Sex = user.Sex
        songList[i].Avatar = user.Avatar
        songList[i].UserId = userid

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
    var vodList []SongMsg = make([]SongMsg,8)
    i := 0

    var rows *sql.Rows
    var result *gorm.DB
    if key == ""{
        result = db.Raw("select id,user_id,name,singer,more,created_at from vod order by rand() limit 30")
        rows,_ = result.Rows()
    }else {
        if sort == "1" {
            result = db.Raw("select id,user_id,name,singer,more,created_at from vod where style=? or language=? order by rand() limit 30",key,key)
            rows,_ = result.Rows()
        }else {
            result = db.Raw("select id,user_id,name,singer,more,created_at from vod where style=? or language=? order by created_at,praise desc limit 30",key,key)
            rows,_ = result.Rows()
        }
    } 

    if rows == nil {
        return vodList
    }

    defer rows.Close()

    for rows.Next() {
        var vod statements.Vod
        db.ScanRows(rows,&vod)
        vodList[i].Id = vod.ID
        vodList[i].Name = vod.Name
        vodList[i].More = vod.More
        vodList[i].Time = vod.CreatedAt
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
    if data1==nil {return MainMsg{},nil}
    if err1!=nil {return MainMsg{},err1}
    data2,err2 := client.Get("healing2020:Main:"+key+"ListenMsg"+sort).Bytes()
    if data2==nil {return MainMsg{},nil}
    if err2!=nil {return MainMsg{},err2}

    var sing []SongMsg
    var listen []SongMsg
    json.Unmarshal(data1,&sing)
    json.Unmarshal(data2,&listen)
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

type SearchResp struct {
    User []UserResp 
    Song []SongResp 
    Vod []VodResp 
    Err error 
}

type UserResp struct {
    UserId uint `json:"userid"`
    UserName string `json:"userName"`
    Avatar string `json:"avatar"`
    More string `json:"more"`
}

type SongResp struct {
    SongId uint `json:"songid"`
    SongName string `json:"songName"`
    Praise int `json:"praise"` 
    Source string `json:"source"`
    Singer string `json:"singer"`
    Time time.Time `json:"time"`
}

type VodResp struct {
    VodId uint `vodid`
    VodName string `json:"vodName"`
    VodUser string `json:"vodUser"`
    Time time.Time `json:"time"`
}

func GetSearchResult(search string) SearchResp{
    db := setting.MysqlConn()
    defer db.Close()

    var searchResp SearchResp
    var songResp []SongResp
    var userResp []UserResp
    var vodResp []VodResp


    var song statements.Song
    //result := db.Raw("select id,praise,source,created_at,user_id from song where name = ?",search)
    result := db.Model(&statements.Song{}).Where("name = ?",search).Select("id,praise,source,created_at,user_id").Find(&song)
    if result.RowsAffected != 0 && result.Error == nil {
        rows,_ := result.Rows()
        defer rows.Close()

        songResp = make([]SongResp,result.RowsAffected)

        i := 0
        for rows.Next() {
            db.ScanRows(rows,&song)
            songResp[i].SongId = song.ID
            songResp[i].Praise = song.Praise
            songResp[i].Source = song.Source
            songResp[i].SongName = search
            songResp[i].Time = song.CreatedAt

            var user statements.User
            db.Model(&statements.User{}).Select("nick_name").Where("id = ?",song.UserId).First(&user)
            songResp[i].Singer = user.NickName

            i++
        }    
    }else {
        searchResp.Err = result.Error
        if errors.Is(result.Error,gorm.ErrRecordNotFound) {
            searchResp.Err = nil
        }
    }

    var vod statements.Vod
    //result = db.Raw("select id,created_at,user_id from vod where name = ?",search)
    result = db.Model(&statements.Vod{}).Where("name = ?",search).Select("id,created_at,user_id").Find(&vod)

    if result.RowsAffected != 0 && result.Error == nil {
        rows,_ := result.Rows()
        defer rows.Close()

        var vodResp []VodResp = make([]VodResp,result.RowsAffected)

        i := 0
        for rows.Next() {
            db.ScanRows(rows,&vod)
            vodResp[i].VodId = vod.ID
            vodResp[i].VodName = search
            vodResp[i].Time = vod.CreatedAt

            var user statements.User
            db.Model(&statements.User{}).Select("nick_name").Where("id = ?",vod.UserId).First(&user)
            vodResp[i].VodUser = user.NickName

            i++
        }    
    }else {
        searchResp.Err = result.Error
        if errors.Is(result.Error,gorm.ErrRecordNotFound) {
            searchResp.Err = nil
        }
    }

    //result = db.Raw("select id,more,avatar from user where nick_name = ?",search)
    var user statements.User
    result = db.Model(&statements.User{}).Where("nick_name=?",search).Select("id,more,avatar").Find(&user)
    if result.RowsAffected != 0 && result.Error == nil {
        rows,_ := result.Rows()
        defer rows.Close()

        userResp = make([]UserResp,result.RowsAffected)

        i := 0
        for rows.Next() {
            db.ScanRows(rows,&user)
            userResp[i].UserId = user.ID
            userResp[i].More = user.More
            userResp[i].Avatar = user.Avatar
            userResp[i].UserName = search

            i++
        }    
    }else {
        searchResp.Err = result.Error
        if errors.Is(result.Error,gorm.ErrRecordNotFound) {
            searchResp.Err = nil
        }
    }

    searchResp.User = userResp
    searchResp.Song = songResp
    searchResp.Vod = vodResp

    return searchResp
}
