package models

import (
	"healing2020/pkg/setting"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Subject struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Intro string `json:"intro"`
}

func SingSubject() ([]Subject, error) {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	//获取歌房信息
	var singSubject []Subject
	err := db.Table("Subject").Select("id, name, intro").Scan(&singSubject).Error
	return singSubject, err
}