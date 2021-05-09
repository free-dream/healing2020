package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"strconv"
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
    IsPraise  bool      `json:"isPraise"`
}

type AllDeliver struct {
	Deliverelse User
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avater"`
}

func DeliverHome(Type string, myID uint) ([]AllDeliver, error) {
	var err error
	//连接mysql
	db := setting.MysqlConn()

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
        responseDeliver[i].Deliverelse.IsPraise,_ = HasPraise(1,myID,uint(deliverHome[i].Id))
        responseDeliver[i].Deliverelse.Praise = GetPraiseCount("deliver",uint(deliverHome[i].Id))
	}

	return responseDeliver, err
}

//发送单个投递详情
func SingleDeliver(DevId string, myID uint) ([]AllDeliver, error) {
	//连接mysql
	db := setting.MysqlConn()

	intId, _ := strconv.Atoi(DevId)
	deliverId := uint(intId)

	var singleDeliver []User
	//获取单个投递信息
	err := db.Table("deliver").Select("id, user_id, created_at, type, text_field, photo, record, praise").Where("id = ? ", deliverId).First(&singleDeliver).Error
	if err != nil {
		return nil, err
	}

	//获取用户昵称
	SingleElse := make([]statements.User, len(singleDeliver))
	for i := 0; i < len(singleDeliver); i++ {
		err2 := db.Table("user").Select("nick_name, avatar").Where("id = ?", singleDeliver[i].UserID).Scan(&SingleElse[i]).Error
		if err2 != nil {
			return nil, err2
		}
	}

	responseSingle := make([]AllDeliver, len(singleDeliver))
	for i := 0; i < len(singleDeliver); i++ {
		responseSingle[i] = AllDeliver{
			Deliverelse: singleDeliver[i],
			Nickname:    SingleElse[i].NickName,
			Avatar:      SingleElse[i].Avatar,
		}
        responseSingle[i].Deliverelse.IsPraise,_ = HasPraise(1,myID,uint(singleDeliver[i].Id))
        responseSingle[i].Deliverelse.Praise = GetPraiseCount("deliver",uint(singleDeliver[i].Id))
	}
	return responseSingle, err
}
