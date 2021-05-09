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

	// status := 0
	tx := db.Begin()
	//投递评论
	if Type == "2" {
		var com statements.Comment
		com.UserId = userid
		com.DeliverId = ID
		com.Type = Typeid
		com.Content = Content
		if err := tx.Model(&statements.Comment{}).Create(&com).Error; err != nil {
			return err
		}
	}

	//歌房评论
	if Type == "1" {
		var com statements.Comment
		com.UserId = userid
		com.SongId = ID
		com.Type = Typeid
		com.Content = Content
		if err := tx.Model(&statements.Comment{}).Create(&com).Error; err != nil {
			return err
		}
	}
	return tx.Commit().Error
}
