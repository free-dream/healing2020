package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"

  _ "github.com/jinzhu/gorm/dialects/mysql"
)

//注册时为用户初始化background
func CreateBackground(userID uint) error {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()
	
	//开启事务
	tx := db.Begin()
	err := tx.Model(&statements.Background{}).Create(statements.Background{UserId: userID}).Error
	
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//更新用户使用的背景
func UpdateBackgroundNow(userID uint, toSaveBackground int) error {
	db := setting.MysqlConn()
	defer db.Close()

	tx := db.Begin()
	err := tx.Model(&statements.Background{}).Where("user_id = ?", userID).Update(statements.Background{Now: toSaveBackground}).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}