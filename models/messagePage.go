package models

import (
	"fmt"
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ToMessagePage struct {
	User     Target `json:"target"`
	LastText Last   `json:"last"`
}

type Target struct {
	ID     uint   `json:"id"`
	Name   string `json:"name" gorm:"column:nick_name"`
	Avatar string `json:"avatar"`
}

type Last struct {
	ID         uint      `json:"id"`
	Content    string    `json:"content"`
	Time       time.Time `json:"time" gorm:"column:created_at"`
	URL        string    `json:"url"`
	StringTime string    `json:"stringtime"`
	Type       int       `json:"type" gorm:"column:type"`
}

//获取和用户聊过天的对象信息
func selectTarget(db *gorm.DB, userID uint) ([]Target, error) {
	//聊天只能由录音发起
	//查录音记录获取和用户聊过天的对象
	var messageSender []statements.Message
	var messageReceive []statements.Message
	err := db.Select("send").Where("receive = ? AND type = ?", userID, 1).Find(&messageSender).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = db.Select("receive").Where("send = ? AND type = ?", userID, 1).Find(&messageReceive).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//将messageReceive的Receive的id数据转到Send中,后面Send就存放着对象的id
	//再将数组合并
	//PS:对value进行的修改不会保存
	for key, value := range messageReceive {
		messageReceive[key].Send = value.Receive
	}
	allTargetID := append(messageSender, messageReceive...)
	//fmt.Println(allTargetID)

	//获取和用户聊过天的对象具体信息
	userTarget := make([]Target, len(allTargetID))
	for key, value := range allTargetID {
		err = db.Table("user").Where("id = ?", value.Send).Scan(&userTarget[key]).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	//fmt.Println(userTarget)
	return userTarget, err
}

//获取最后一条信息
func selectLast(db *gorm.DB, userTarget []Target, userID uint) ([]Last, error) {
	//分别获取用户与聊天对象各自发出的最新消息
	//比对时间确定最后一条消息
	var sendMessage Last
	var receiveMessage Last
	lastMessage := make([]Last, len(userTarget))
	var err error

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

		sendMessage.StringTime = sendMessage.Time.Format("2006-01-02 15:04:05")
		receiveMessage.StringTime = receiveMessage.Time.Format("2006-01-02 15:04:05")

		if sendMessage.Time.After(receiveMessage.Time) {
			lastMessage[key] = sendMessage
		} else {
			lastMessage[key] = receiveMessage
		}
	}
	return lastMessage, err
}

//综合并返回消息首页信息
func ResponseMessagePage(userID uint) ([]ToMessagePage, error) {
	db := setting.MysqlConn()
	defer db.Close()

	userTarget, err := selectTarget(db, userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	lastMessage, err := selectLast(db, userTarget, userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
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
