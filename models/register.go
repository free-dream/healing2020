package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"

  _ "github.com/jinzhu/gorm/dialects/mysql"
)


func RegisterUpdate(user statements.User) error {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()
	//开启事务
	tx := db.Begin()
	err := tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}