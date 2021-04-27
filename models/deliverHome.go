package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	UserID    int    `json:"user_id" `
	Type      int    `json:"Type"`
	TextField string `json:"text_field"`
	Photo     string `json:"photo"`
	Record    string `json:"record"`
	Praise    int    `json:"praise"`
}

type AllDeliver struct {
	Deliverelse User
	Nickname    string `json:"nickname"`
}

func DeliverHome(Type string) ([]AllDeliver, error) {
	var err error
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	var deliverHome []User
	//最新排序
	if Type == "0" {
		//获取投递信息
		err := db.Table("deliver").Select("user_id, created_at, type, text_field, photo, record, praise").Order("created_at DESC").Scan(&deliverHome).Error
		if err != nil {
			return nil, err
		}
	}
	//随机排序
	if Type == "1" {
		err := db.Table("deliver").Select("user_id, type, text_field, photo, record, praise").Order("rand()").Scan(&deliverHome).Error
		if err != nil {
			return nil, err
		}
	}

	//获取用户昵称
	UserNickname := make([]statements.User, len(deliverHome))
	for i := 0; i < len(deliverHome); i++ {
		err = db.Table("user").Select("nick_name").Where("id = ?", deliverHome[i].UserID).Scan(&UserNickname[i]).Error
		if err != nil {
			return nil, err
		}
	}

	responseDeliver := make([]AllDeliver, len(deliverHome))
	for i := 0; i < len(deliverHome); i++ {
		responseDeliver[i] = AllDeliver{
			Deliverelse: deliverHome[i],
			Nickname:    UserNickname[i].NickName,
		}
	}

	return responseDeliver, err
}
