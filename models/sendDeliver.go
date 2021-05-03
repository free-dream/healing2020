package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func PostDeliver(UserId string, TextField string, Photo string, Record string) error {
	intId, _ := strconv.Atoi(UserId)
	userid := uint(intId)
	db := setting.MysqlConn()
	defer db.Close()

	// status := 0
	tx := db.Begin()

	if TextField != "" {
		//发送纯文字投递
		if Photo == "" && Record == "" {
			var dev statements.Deliver
			dev.UserId = userid
			dev.TextField = TextField
			dev.Type = 1
			if err := tx.Model(&statements.Deliver{}).Create(&dev).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		//发送图文投递
		if Photo != "" && Record == "" {
			var dev statements.Deliver
			dev.UserId = userid
			dev.TextField = TextField
			dev.Photo = Photo
			dev.Type = 2
			if err := tx.Model(&statements.Deliver{}).Create(&dev).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		//发送带录音投递
		if Record != "" {
			var dev statements.Deliver
			dev.UserId = userid
			dev.TextField = TextField
			dev.Type = 3
			dev.Photo = Photo
			dev.Record = Record
			if err := tx.Model(&statements.Deliver{}).Create(&dev).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}
