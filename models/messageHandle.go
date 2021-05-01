package models

import (
	"fmt"
	"healing2020/models/statements"
	"healing2020/pkg/setting"
)

//保存用户聊天的消息
func CreateMessage(msg statements.Message) error {
	db := setting.MysqlConn()
	defer db.Close()

	tx := db.Begin()
	err := tx.Model(&statements.Message{}).Create(msg).Error
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return err
	}
	return tx.Commit().Error
}
