package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"strconv"

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

//发送歌房数据
func PostSubject(ID string, Name string, Photo string, Intro string) error {
	intId, _ := strconv.Atoi(ID)
	subject_id := uint(intId)
	db := setting.MysqlConn()

	status := 0
	tx := db.Begin()
	if Name != "" {
		//发送
		var dev statements.Subject
		dev.ID = subject_id
		dev.Name = Name
		dev.Photo = Photo
		dev.Intro = Intro
		err := tx.Model(&statements.Subject{}).Create(&dev).Error
		if err != nil {
			if status < 5 {
				status++
				tx.Rollback()
			} else {
				return err
			}
		}
	}
	return tx.Commit().Error
}
