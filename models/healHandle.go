package models

import (
    "strconv"
    "healing2020/pkg/setting"
    "healing2020/pkg/tools"
    "healing2020/models/statements"
)

func GetPhone(info tools.RedisUser) string{
    return info.Phone
}

type RecordResp struct {
    Source string `json:"source"`
    Err error `json:"err"`
}
func GetRecord(id string) RecordResp{
    intId,_ := strconv.Atoi(id) 
    songId := uint(intId)
    var song statements.Song
    
    db := setting.MysqlConn()
    defer db.Close()

    result := db.Model(&statements.Song{}).Where("ID=?",songId).First(&song)
    var recordResp RecordResp
    recordResp.Source = song.Source 
    recordResp.Err = result.Error
    return recordResp
}

func CreateRecord(id string,source string,uid uint) error{
    intId,_ := strconv.Atoi(id) 
    vodId := uint(intId)
    db := setting.MysqlConn()
    defer db.Close()
    userId := uid
    status := 0

    tx := db.Begin()

    var vod statements.Vod
    result1 := tx.Model(&statements.Vod{}).Where("ID=?",vodId).First(&vod)
    if result1.Error != nil {
        return nil
    }
    var song statements.Song
    song.VodId = vodId
    song.UserId = userId
    song.VodSend = vod.UserId
    song.Name = vod.Name
    song.Praise = 0
    song.Source = source
    song.Style = vod.Style
    song.Language = vod.Language

    err := tx.Model(&statements.Song{}).Create(&song).Error
    if err != nil {
        if status < 5 {
            status++
            tx.Rollback()
        }else {
            return err
        }
    }
    return tx.Commit().Error
}

func AddPraise(strId string,types string) error{
    intId,_ := strconv.Atoi(strId) 
    id := uint(intId)
    db := setting.MysqlConn()
    defer db.Close()

    tx := db.Begin()
    status := 0
    if types == "1" {
        var song statements.Song
        result := tx.Model(&statements.Song{}).Where("ID=?",id).First(&song)
        if result.Error != nil {
            return result.Error
        }
        song.Praise = song.Praise + 1
        err := tx.Save(&song).Error
        if err != nil {
            if status < 5 {
                status++
                tx.Rollback()
            }else {
                return err
            }
        }
    }
    if types == "2" {
        var deliver statements.Deliver
        result := tx.Model(&statements.Deliver{}).Where("ID=?",id).First(&deliver)
        if result.Error != nil {
            return result.Error
        }
        deliver.Praise = deliver.Praise + 1
        err := tx.Save(&deliver).Error
        if err != nil {
            if status < 5 {
                status++
                tx.Rollback()
            }else {
                return err
            }
        }
    }
    return tx.Commit().Error
}

func CreateVod(uid uint,singer string,style string,language string,name string,more string) error{
    db := setting.MysqlConn()
    defer db.Close()

    status := 0
    tx := db.Begin()
    var vod statements.Vod
    vod.UserId = uid   
    vod.More = more
    vod.Name = name
    vod.Singer = singer
    vod.Style = style
    vod.Language = language
    err := tx.Model(&statements.Vod{}).Create(&vod).Error
    if err != nil {
        if status < 5 {
            status++
            tx.Rollback()
        }else {
            return err
        }
    }
    return tx.Commit().Error
}
