package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"

	//"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func GodAddPrize(name string, photo string, intro string, weight int, count int) error {

	db := setting.MysqlConn()

	tx := db.Begin()
	//发送
	var dev statements.Prize
	dev.Name = name
	dev.Photo = photo
	dev.Intro = intro
	dev.Weight = weight
	dev.Count = count

	err := tx.Model(&statements.Prize{}).Create(&dev).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
