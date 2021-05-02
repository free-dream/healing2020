package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"log"
)

//保存用户聊天的消息
func SaveMessage(msg statements.Message) error {
	db := setting.MysqlConn()
	defer db.Close()

	tx := db.Begin()
	err := tx.Model(&statements.Message{}).Create(&msg).Error
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	}
	return tx.Commit().Error
}

//删除用户聊天消息
func DeleteMessage(msg statements.Message) error {
	db := setting.MysqlConn()
	defer db.Close()

	tx := db.Begin()
	err := tx.Model(&statements.Message{}).Where("msg_id=?", msg.MsgID).Delete(&statements.Message{}).Error
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	}
	return tx.Commit().Error
}
