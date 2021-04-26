package models

import (
	"healing2020/models/statements"
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

type AllDeliver struct {
	Deliverelse User
	Nickname  string    `json:"nickname"`
}
func DeliverHome() ([]AllDeliver, error) {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	//获取投递信息
	var deliverHome []User
	err := db.Table("deliver").Select("user_id, type, text_field, photo, record, praise").Scan(&deliverHome).Error

	//获取用户昵称
	UserNickname := make([]statements.User, len(deliverHome))
	for i := 0; i < len(deliverHome); i++ {
		err = db.Table("User").Select("nick_name").Where("id = ?", deliverHome[i].UserID).Scan(&UserNickname[i]).Error
		if err != nil {
			return nil, err
		}
	}

	responseDeliver := make([]AllDeliver, len(deliverHome))
	for i := 0; i < len(deliverHome); i++ {
		responseDeliver[i] = AllDeliver{
			Deliverelse:    deliverHome[i],
			Nickname:  		UserNickname[i].NickName,
		}
	}
	return responseDeliver, err

}
