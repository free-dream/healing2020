package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type AllSpecial struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Intro    string `json:"intro"`
	Photo	 string `json:"photo"`
	SingSong []UserMessage
}

type Sing struct {
	UserID		uint   		`json:"user_id"`
	Id 			int			`json:"id"`
	CreatedAt	time.Time 	`json:"createdat"`
	Praise 		string 		`json:"praise"`
	Song   		string 		`json:"name"`
}

type UserMessage struct {
	Nickname	string `json:"user"`
	Avatar		string `json:"avatar"`
	SingElse	Sing
	
	//无用数据
	Type 		string `json:"type"`
	Photo 		string `json:"photo"`
	Text		string `json:"text"`
	Source		string `json:"source"`
}
func SingHome(subjectID uint) (AllSpecial, error) {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	//获取歌房信息
	var singSubject statements.Subject
	err := db.Table("subject").Select("id, name, intro, photo").Where("id = ? ", subjectID).First(&singSubject).Error

	var SingHome []Sing
	err = db.Table("special").Select("user_id, id, created_at, praise, song").Where("subject_id = ? ", subjectID).Scan(&SingHome).Error

	UserElse := make([]statements.User, len(SingHome))
	for i := 0; i < len(SingHome); i++ {
		err = db.Table("user").Select("nick_name, avatar").Where("id = ?", SingHome[i].UserID).Scan(&UserElse[i]).Error
	}
	
	responseSing := make([]UserMessage, len(SingHome))
	for i := 0; i < len(SingHome); i++ {
		responseSing[i] = UserMessage{
			Nickname:   UserElse[i].NickName,
			Avatar:     UserElse[i].Avatar,
			SingElse: 	SingHome[i],
		}
	}

	allSpecial := AllSpecial{
		ID:       singSubject.ID,
		Name:     singSubject.Name,
		Intro:    singSubject.Intro,
		SingSong: responseSing,
	}

	return allSpecial, err
}
