package models

import (
	"fmt"
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Money struct {
	Money int `json:"money"`
}

type Task struct {
	Lo1 int `json:"login"`
	Lo2 int `json:"chooseSong"`
	Lo3 int `json:"healing"`
	Lo4 int `json:"singHome"`
	Lo5 int `json:"praise"`
	Lo6 int `json:"share"`
}

//查询当前积分
func GetMoney(userID uint) ([]Money, error) {
	//连接mysql
	db := setting.MysqlConn()

	//获取个人积分信息
	var user []Money
	err := db.Table("user").Select("money").Where("id= ? ", userID).First(&user).Error
	return user, err
}

//抽奖--减少当前积分
func UseMoney(userID uint) error {
	//连接mysql
	db := setting.MysqlConn()
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

//返回任务列表
func GetTask(userID uint) ([]Task, error) {
	//连接mysql
	db := setting.MysqlConn()

	//获取个人积分信息
	var user []Task
	err := db.Table("user_other").Select("lo1, lo2, lo3, lo4, lo5, lo6").Where("user_id= ? ", userID).First(&user).Error
	return user, err
}

//更新每日任务
func UpdateTask() error {
	//连接mysql
	db := setting.MysqlConn()

	//更新每日任务
	err := db.Table("user_other").Update(map[string]interface{}{"lo1": "0", "lo2": "0", "lo3": "0", "lo4": "0", "lo5": "0", "lo6": "0"}).Error
	return err
}

//分享二维码加积分
func PostQRcode(User_id string) error {
	intId, _ := strconv.Atoi(User_id)
	user_id := uint(intId)

	db := setting.MysqlConn()

	status := 0
	tx := db.Begin()

	var userother statements.UserOther
	result := tx.Model(&statements.UserOther{}).Where("user_id= ?", user_id).First(&userother)
	if result.Error != nil {
		return result.Error
	}
	//判断完成每日任务和增加积分
	if userother.Lo6 != 1 {
		err2 := tx.Model(&statements.UserOther{}).Where("user_id = ?", user_id).Update("lo6", 1).Error
		if err2 != nil {
			if status < 5 {
				status++
				tx.Rollback()
			} else {
				return err2
			}
		}
		var user statements.User
		result := tx.Model(&statements.User{}).Where("id= ?", user_id).First(&user)
		if result.Error != nil {
			return result.Error
		}
		if user.Money >= 0 {
			user.Money = user.Money + 10
			err3 := tx.Save(&user).Error
			if err3 != nil {
				if status < 5 {
					status++
					tx.Rollback()
				} else {
					return err3
				}
			}
		}
	}
	return tx.Commit().Error
}
