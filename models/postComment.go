package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func PostComment(UserId string, Id string, Type string, Content string) error {
	intId, _ := strconv.Atoi(UserId)
	userid := uint(intId)
	_Id, _ := strconv.Atoi(Id)
	ID := uint(_Id)
	TypeId, _ := strconv.Atoi(Type)
	Typeid := int(TypeId)
	db := setting.MysqlConn()
	defer db.Close()

	status := 0
	tx := db.Begin()
	//投递评论
	if Type == "2" {
		var com statements.Comment
		com.UserId = userid
		com.DeliverId = ID
		com.Type = Typeid
		com.Content = Content
		err := tx.Model(&statements.Comment{}).Create(&com).Error
		if err != nil {
			if status < 5 {
				status++
				tx.Rollback()
			} else {
				return err
			}
		}
	}

	//歌房评论
	if Type == "1" {
		var com statements.Comment
		com.UserId = userid
		com.SongId = ID
		com.Type = Typeid
		com.Content = Content
		err := tx.Model(&statements.Comment{}).Create(&com).Error
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
