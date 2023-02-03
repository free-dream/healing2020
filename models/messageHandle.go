package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
)

//保存用户聊天的消息
func SaveMessage(msg statements.Message) error {
	db := setting.MysqlConn()

	err := db.Model(&statements.Message{}).Create(&msg).Error
	if err != nil {
		return err
	}
	return nil
}

//删除用户聊天消息
func DeleteMessage(msg statements.Message) error {
	db := setting.MysqlConn()

	err := db.Model(&statements.Message{}).Where("msg_id=? AND is_to_from_user_id=?", msg.MsgID, msg.IsToFromUserID).Delete(&statements.Message{}).Error
	if err != nil {
		return err
	}
	return nil
}

//获取message里的所有信息
func SelectAllMessage() ([]statements.Message, error) {
	db := setting.MysqlConn()

	var allMessage []statements.Message
	err := db.Find(&allMessage).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return allMessage, nil
}
