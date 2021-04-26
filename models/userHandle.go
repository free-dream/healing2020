package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"errors"
	//"encoding/json"
	//"time"
)

func UpdateOrCreate(openId string, nickName string, sex int, avatar string) {
	db := setting.MysqlConn()
	db.Transaction(func(tx *gorm.DB) error {
		var user statements.User
		result := tx.Model(&statements.User{}).Where("open_id=?", openId).First(&user)
		user.NickName = nickName
		user.OpenId = openId
		user.Avatar = avatar
		user.Sex = sex
		var result2 *gorm.DB
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			result2 = tx.Model(&statements.User{}).Create(&user)
		} else {
			result2 = tx.Model(&statements.User{}).Where("open_id=?", openId).Update(&user)
		}
		// client := setting.RedisConn()
		// defer client.Close()
		// dataByte,_ := json.Marshal(user)
		// data := string(dataByte)
		// keyname := "healing2020:token:"+token
		// client.Set(keyname,data,time.Minute*30)

		return result2.Error
	})
}

//SELECT hobby FROM user where id = userID
func SelectUserHobby(userID uint) (string, error) {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	var userHobby statements.User
	err := db.Select("hobby").Where("id=?", userID).First(&userHobby).Error
	return userHobby.Hobby, err
}

//更新user表
func UpdateUser(user statements.User, userID uint) error {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	//开启事务
	tx := db.Begin()
	err := tx.Model(&statements.User{}).Where("id=?", userID).Update(user).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
