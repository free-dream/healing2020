package models

import (
	"healing2020/pkg/setting"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	UserID    int    `json:"user_id" `
	Type      int    `json:"Type"`
	Textfield string `json:"textfield"`
	Photo     string `json:"photo"`
	Record    string `json:"record"`
	Praise    int    `json:"praise"`
}

func DeliverHome() ([]User, error) {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	//获取投递信息
	var deliverHome []User
	err := db.Table("deliver").Select("user_id, type, text_field, photo, record, praise").Scan(&deliverHome).Error
	return deliverHome, err
}
