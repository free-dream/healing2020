package models

import (
	"fmt"
	"healing2020/models/statements"
	"healing2020/pkg/setting"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Money struct {
	Money int `json:"money"`
}

//查询当前积分
func GetMoney(userID uint) ([]Money, error) {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	//获取个人积分信息
	var user []Money
	err := db.Table("user").Select("money").Where("id= ? ", userID).First(&user).Error
	return user, err
}

//抽奖--减少当前积分
func UseMoney(userID uint) error {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()
	//进行抽奖
	status := 0
	tx := db.Begin()
	var user statements.User
	result := tx.Model(&statements.User{}).Where("id= ?", userID).First(&user)
	if result.Error != nil {
		return result.Error
	}
	if user.Money >= 100 {
		user.Money = user.Money - 100
		err := tx.Save(&user).Error
		if err != nil {
			if status < 5 {
				status++
				tx.Rollback()
			} else {
				return err
			}
		}
	} else {
		err := fmt.Errorf("")
		return err
	}
	return tx.Commit().Error
}

//每日任务--增加当前积分
func EarnMoney(userID uint) error {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()
	//每日任务获取积分
	status := 0
	tx := db.Begin()
	var user statements.User
	result := tx.Model(&statements.User{}).Where("id= ?", userID).First(&user)
	if result.Error != nil {
		return result.Error
	}
	if user.Money >= 0 {
		user.Money = user.Money + 30
		err := tx.Save(&user).Error
		if err != nil {
			if status < 5 {
				status++
				tx.Rollback()
			} else {
				return err
			}
		}
	} else {
		err := fmt.Errorf("")
		return err
	}
	return tx.Commit().Error
}