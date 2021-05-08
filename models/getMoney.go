package models

import (
	"fmt"
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"log"
	"math/rand"
	"strconv"
	"time"

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

type GetPrize struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
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
	result := tx.Model(&statements.UserOther{}).Where("user_id = ?", user_id).First(&userother)
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

//真·抽奖
func LotteryDraw(userID uint, bd []int, bdstr string) (GetPrize, error) {
	//连接mysql
	db := setting.MysqlConn()
	var responseLottery GetPrize
	//进行抽奖
	tx := db.Begin()

	rand.Seed(time.Now().UnixNano())
	prop := rand.Intn(999)
	log.Println(prop)

	var black GetPrize

	var prize []Prize
	err := tx.Model(&statements.Prize{}).Select("id, name, photo, weight, count").Scan(&prize).Error
	if err != nil {
		log.Println(err)
	}

	var user statements.User
	result0 := tx.Model(&statements.User{}).Where("id= ?", userID).First(&user)
	if result0.Error != nil {
		return black, result0.Error
	}

	log.Println(user.Money)

	if user.Money >= 100 {
	LOOP:
		for i := 0; i < len(prize); i++ {
			if prop <= prize[i].Weight && prize[i].Count > 0 {
				responseLottery = GetPrize{
					Id:    prize[i].Id,
					Name:  prize[i].Name,
					Photo: prize[i].Photo,
				}

				//增加匿名次数
				if prize[i].Id == 6 {
					var anonymous statements.UserOther
					result1 := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).First(&anonymous)
					if result1.Error != nil {
						return black, result1.Error
					}
					anonymous.RemainHideName = anonymous.RemainHideName + 1
					err0 := tx.Save(&anonymous).Error
					if err0 != nil {
						tx.Rollback()
					}
				}

				//增加背景图
				if prize[i].Id == 5 {
					bdcount := len(bd)
					if bdcount == 6 {
						i = i + 1
						responseLottery = GetPrize{
							Id:    prize[i].Id,
							Name:  prize[i].Name,
							Photo: prize[i].Photo,
						}		
					} else {
						bdstr = bdstr + "," + strconv.Itoa(len(bd) + 1)
						result2 := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).Update("ava_background", bdstr).Error
						if result2 != nil {
							tx.Rollback()
						}
					}
				}
				
				//修改剩余数量
				var count statements.Prize
				result3 := tx.Model(&statements.Prize{}).Where("id = ?", uint(prize[i].Id)).First(&count)
				if result3.Error != nil {
					return black, result3.Error
				}
				count.Count = count.Count - 1
				err1 := tx.Save(&count).Error
				if err1 != nil {
					tx.Rollback()
				}

				//存入我的奖品
				var lot statements.Lottery
				lot.PrizeId = uint(prize[i].Id)
				lot.UserId = userID
				if result4 := tx.Model(&statements.Lottery{}).Create(&lot).Error; err != nil {
					tx.Rollback()
					return black, result4
				}

				//修改剩余积分
				user.Money = user.Money - 100
				err2 := tx.Save(&user).Error
				if err2 != nil {
					tx.Rollback()
				}
				break LOOP
			}
		}
	} else {
		log.Println(result0)
	}
	return responseLottery, tx.Commit().Error
}
