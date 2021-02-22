package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"

  _ "github.com/jinzhu/gorm/dialects/mysql"
)

//更新user表
func UpdateUser(user statements.User, userID uint) error{
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()
	//开启事务
	tx := db.Begin()
	err := tx.Model(&statements.User{}).Where("id=?",userID).Update(user).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}