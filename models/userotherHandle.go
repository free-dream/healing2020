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

//获取userother
func SelectUseOther(userID uint) (statements.UserOther, error) {
	//连接mysql
	db := setting.MysqlConn()

	var userOther statements.UserOther
	err := db.Where("user_id=?", userID).First(&userOther).Error
	return userOther, err
}
