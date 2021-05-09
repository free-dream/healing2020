package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//更新用户使用的背景
func UpdateUserOtherNow(userID uint, toSaveUserOther int) error {
	db := setting.MysqlConn()

	tx := db.Begin()
	err := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).Update(statements.UserOther{Now: toSaveUserOther}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

//获取剩余点歌和匿名次数
func SelectRemainNum(userID uint) (statements.UserOther, error) {
	//连接mysql
	db := setting.MysqlConn()
	var userOther statements.UserOther
	err := db.Select("remain_sing, remain_hide_name").Where("user_id=?", userID).First(&userOther).Error
	return userOther, err
}

//每日更新用户剩余点歌次数
func UpdateRemainSingDay() {
	db := setting.MysqlConn()
	db.Model(&statements.UserOther{}).Where("remain_sing < ?", 8).Update("remain_sing", 8)
}
