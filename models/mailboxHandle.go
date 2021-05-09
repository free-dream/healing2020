package models

import (
	"fmt"
	"healing2020/models/statements"
	"healing2020/pkg/setting"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//保存广播信息
func CreateMailBox(message string) error {
	db := setting.MysqlConn()

	tx := db.Begin()
	err := tx.Model(&statements.Mailbox{}).Create(&statements.Mailbox{Message: message}).Error
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
