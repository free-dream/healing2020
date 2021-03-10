package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type AllSpecial struct {
	ID    uint    `json:"id"`
	Name  string `json:"name"`
	Intro string `json:"intro"`	
	SingSong []Sing
}

type Sing struct {
	UserID    uint    `json:"id"`
	Name  string `json:"name"`
	Praise string `json:"praise"`
	Song string `json:"song"`
}

func SingHome(subjectID uint) (AllSpecial, error) {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	//获取歌房信息
	var singSubject statements.Subject
	err := db.Table("Subject").Select("id, name, intro").Where("id = ? ", subjectID).First(&singSubject).Error

	var SingHome []Sing
	err = db.Table("Special").Select("user_id, name, praise, song").Where("subject_id = ? ", subjectID).Scan(&SingHome).Error
	
	// var allSpecial []AllSpecial

	// 	allSpecial.ID = singSubject.id
	// 	allSpecial.Name = singSubject.name
	// 	allSpecial.Intro = singSubject.Intro
	// 	allSpecial.SingSong = SingHome[3]
	allSpecial := AllSpecial{
		ID: singSubject.ID,
		Name: singSubject.Name, 
		Intro: singSubject.Intro, 
		SingSong: SingHome,
	}

	return allSpecial, err
}
