package models

import (
	"fmt"
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ToMessagePage struct {
	User     Target `json:"target"`
	LastText Last   `json:"last"`
}

type Target struct {
	ID     uint   `json:"id"`
	Name   string `gorm:"column:nick_name"`
	Avatar string
}

type Last struct {
	ID       uint `json:"id"`
	Content  string
	Time     time.Time `gorm:"column:created_at"`
	TextType int       `json:"type" gorm:"column:type"`
}

func ResponseMessagePage(userID uint) ([]ToMessagePage, error) {
	db := setting.MysqlConn()
	defer db.Close()

	//聊天只能由录音发起
	//查录音记录获取和用户聊过天的对象
	var messageSender []statements.Message
	var messageReceive []statements.Message
	err := db.Select("send").Where("receive = ? AND type = ?", userID, 1).Find(&messageSender).
		Select("receive").Where("send = ? AND type = ?", userID, 1).Find(&messageReceive).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	allTargetID := append(messageSender, messageReceive...)

	//获取和用户聊过天的对象具体信息
	userTarget := make([]Target, len(allTargetID))
	for key, value := range allTargetID {
		err = db.Table("user").Select("id, nick_name, avatar").Where("id = ?", value.Send).Scan(&userTarget[key]).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	fmt.Println(userTarget)

	//分别获取用户与聊天对象各自发出的最新消息
	//比对时间确定最后一条消息
	var sendMessage Last
	var receiveMessage Last
	lastMessage := make([]Last, len(userTarget))
	for key, value := range userTarget {
		err = db.Table("message").Where("send = ? AND receive = ?", userID, value.ID).Order("created_at desc").Limit(1).Scan(&sendMessage).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		err = db.Table("message").Where("send = ? AND receive = ?", value.ID, userID).Order("created_at desc").Limit(1).Scan(&receiveMessage).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		fmt.Println(sendMessage)
		fmt.Println(receiveMessage)
		if sendMessage.Time.After(receiveMessage.Time) {
			lastMessage[key] = sendMessage
		} else {
			lastMessage[key] = receiveMessage
		}
	}

	//综合数据return
	responseMessage := make([]ToMessagePage, len(userTarget))
	for i := 0; i < len(userTarget); i++ {
		responseMessage[i] = ToMessagePage{
			User:     userTarget[i],
			LastText: lastMessage[i],
		}
	}
	return responseMessage, err
}
