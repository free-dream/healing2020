package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	Id        int       `json:"deliver_id"`
	UserID    int       `json:"user_id" `
	CreatedAt time.Time `json:"created_at"`
	Type      int       `json:"Type"`
	TextField string    `json:"text_field"`
	Photo     string    `json:"photo"`
	Record    string    `json:"record"`
	Praise    int       `json:"praise"`
}

type AllDeliver struct {
	Deliverelse User
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avater"`
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
		err := db.Table("deliver").Select("id, user_id, created_at, type, text_field, photo, record, praise").Order("created_at DESC").Scan(&deliverHome).Error
		if err != nil {
			return nil, err
		}
	}
	//随机排序
	if Type == "1" {
		err := db.Table("deliver").Select("id, user_id, created_at, type, text_field, photo, record, praise").Order("rand()").Scan(&deliverHome).Error
		if err != nil {
			return nil, err
		}
	}

	//获取用户昵称
	UserElse := make([]statements.User, len(deliverHome))
	for i := 0; i < len(deliverHome); i++ {
		err = db.Table("user").Select("nick_name, avatar").Where("id = ?", deliverHome[i].UserID).Scan(&UserElse[i]).Error
		if err != nil {
			return nil, err
		}
	}

	responseDeliver := make([]AllDeliver, len(deliverHome))
	for i := 0; i < len(deliverHome); i++ {
		responseDeliver[i] = AllDeliver{
			Deliverelse: deliverHome[i],
			Nickname:    UserElse[i].NickName,
			Avatar:      UserElse[i].Avatar,
		}
	}

	return responseDeliver, err
}
