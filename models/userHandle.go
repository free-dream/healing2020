package models

import (
    "healing2020/models/statements"
    "healing2020/pkg/setting"
    
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"

    "errors"
)

func UpdateOrCreate(openId string,nickName string) {
    db := setting.MysqlConn()
    db.Transaction(func(tx *gorm.DB) error {
        var user statements.User
        result := tx.Model(&statements.User{}).Where("open_id=?",openid).First(&user)
        user.NickName = nickName
        user.OpenId = openId
        var result2 *gorm.DB
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            result2 = tx.Model(&statements.User{}).Create(&user)
        }else {
            result2 = tx.Model(&statements.User{}).Where("open_id=?",openId).Update(&user)
        }
        return result2.Error
    })
}
