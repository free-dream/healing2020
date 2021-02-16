package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"

  _ "github.com/jinzhu/gorm/dialects/mysql"
)

func HobbyUpdate(hobby string, userID uint) error {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()
	//开启事务
	tx := db.Begin()
	err := tx.Model(&statements.User{}).Where("id = ?", userID).Update("Hoppy", hobby).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}