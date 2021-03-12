package models

import (
	"fmt"
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"time"
)

//聊天室具体消息接口返回数据
type ToMessageCell struct {
	ID         uint      `json:"id"`
	ToUserID   uint      `json:"toUserID" gorm:"column:send"`
	FromUserID uint      `json:"fromUserID gorm:"column:receive"`
	Content    string    `json:"content"`
	Time       time.Time `json:"time" gorm:"column:created_at"`
	StringTime string    `json:"stringtime"`
	URL        string    `json:"url"`
	Type       int       `json:"type" gorm:"column:type"`
}

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

//返回聊天室的具体消息
func SelectCellMessage(userID uint, targetID uint) ([]ToMessageCell, error) {
	db := setting.MysqlConn()
	defer db.Close()

	var msgUserSend []ToMessageCell
	var msgTargetSend []ToMessageCell
	err := db.Table("message").Where("send = ? AND receive = ?", userID, targetID).Scan(&msgUserSend).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = db.Table("message").Where("send = ? AND receive = ?", targetID, userID).Scan(&msgTargetSend).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	msg := append(msgUserSend, msgTargetSend...)

	for key, value := range msg {
		msg[key].StringTime = value.Time.Format("2006-01-02 15:04:05")
	}

	return msg, err

}
