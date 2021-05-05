package models

import (
	"healing2020/pkg/setting"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Subject struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Intro string `json:"intro"`
	Photo string `json:"photo"`
}

func SingSubject() ([]Subject, error) {
	//连接mysql
	db := setting.MysqlConn()

	//获取歌房信息
	var singSubject []Subject
	err := db.Table("subject").Select("id, name, intro, photo").Scan(&singSubject).Error
	return singSubject, err
}
