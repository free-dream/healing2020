package models

import (
    "healing2020/pkg/setting"
    "healing2020/models/statements"
)

type FinishData struct {
    Days int `json:"day"`
    AllUser int `json:"allUser"`
    AllSong int `json:"allSong"`

    MostSong string `json:"mostSong"`
    MostVod string `json:"mostVod"`

    HotTags []string `json:"hotTags"`
    HotSong string `json:"hotSong"`
    HotUser string `json:"hotUser"`

    MyVod int `json:"myVod"`
    MyGotHeal int `json:"myGotHeal"`
    MyHeal int `json:"myHeal"`
    MyPraise int `json:"myPraise"`
    MyGotPraise int `json:"myGotPraise"`
}
func GetFinish(userid uint) FinishData{
    var result FinishData
    result.Days = 14
    result.HotTags = []string{""}
    db := setting.MysqlConn()
    count := 0
    db.Model(&statements.User{}).Count(&count)
    result.AllUser = count
    db.Model(&statements.Song{}).Count(&count)
    result.AllSong = count

    var song statements.Song
    db.Raw("select name from song group by name order by count(*) desc").Limit(1).Scan(&song)
    result.MostSong = song.Name
    var vod statements.Vod
    db.Raw("select name from vod group by name order by count(*) desc").Limit(1).Scan(&vod)
    result.MostVod = vod.Name
    db.Model(&statements.Song{}).Select("name").Order("praise desc").First(&song)
    result.HotSong = song.Name 
    var user statements.User
    db.Model(&statements.User{}).Select("nick_name").Order("money desc").First(&user)
    result.HotUser = user.NickName

    db.Model(&statements.Vod{}).Where("user_id = ?", userid).Count(&count)
    result.MyVod = count
    db.Model(&statements.Song{}).Where("vod_send = ?",userid).Count(&count)
    result.MyGotHeal = count
    db.Model(&statements.Song{}).Where("user_id = ?",userid).Count(&count)
    result.MyHeal = count
    db.Model(&statements.Praise{}).Where("type < 4 and is_cancel = 0 and user_id = ?",userid).Count(&count)
    result.MyPraise = count
    db.Raw("select * from praise where is_cancel = 0 and type = 2 and (select user_id from song where song.id = praise.praise_id) = ?", userid).Count(&count)
    result.MyGotPraise = count

    return result
}
