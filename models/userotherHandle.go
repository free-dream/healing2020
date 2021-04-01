package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//注册时为用户初始化杂项信息
func CreateUserOther(userID uint) error {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	//开启事务
	tx := db.Begin()
	err := tx.Model(&statements.UserOther{}).Create(statements.UserOther{UserId: userID}).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//更新用户使用的背景
func UpdateUserOtherNow(userID uint, toSaveUserOther int) error {
	db := setting.MysqlConn()
	defer db.Close()

	tx := db.Begin()
	err := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).Update(statements.UserOther{Now: toSaveUserOther}).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
