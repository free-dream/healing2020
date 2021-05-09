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

//完成每日任务
func FinishTask(task string, userID uint) error {
	//连接mysql
	db := setting.MysqlConn()
	tx := db.Begin()

	switch task {
	case "1" :
		//login
		var userother statements.UserOther
		result := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).First(&userother)
		if result.Error != nil {
			return result.Error
		}
		if userother.Lo1 != 1 {
			t := time.Now()
			t_zero := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
			t_to_tomorrow := 24*60*60 - (t.Unix() - t_zero)		
			redis_cli := setting.RedisClient
			logined := !redis_cli.SetNX(fmt.Sprintf("finish_lo1_id:%d", userID), 0, time.Duration(t_to_tomorrow)*time.Second).Val()
			
			if !logined {
				err2 := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).Update("lo1", 1).Error
				if err2 != nil {
					tx.Rollback()
					return err2
				}
				
				var user statements.User
				result := tx.Model(&statements.User{}).Where("id= ?", userID).First(&user)
				if result.Error != nil {
					return result.Error
				}
				if user.Money >= 0 {
					user.Money = user.Money + 50
					err3 := tx.Save(&user).Error
					if err3 != nil {
						tx.Rollback()
						return err3
					}
				}
			}
		}
		break
	case "2" :
		//vod
		var userother statements.UserOther
		result := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).First(&userother)
		if result.Error != nil {
			return result.Error
		}
		if userother.Lo2 != 1 {
			t := time.Now()
			t_zero := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
			t_to_tomorrow := 24*60*60 - (t.Unix() - t_zero)		
			redis_cli := setting.RedisClient
			voded := !redis_cli.SetNX(fmt.Sprintf("finish_lo2_id:%d", userID), 0, time.Duration(t_to_tomorrow)*time.Second).Val()
			
			if !voded {
				err2 := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).Update("lo2", 1).Error
				if err2 != nil {
					tx.Rollback()
					return err2
				}
				
				var user statements.User
				result := tx.Model(&statements.User{}).Where("id= ?", userID).First(&user)
				if result.Error != nil {
					return result.Error
				}
				if user.Money >= 0 {
					user.Money = user.Money + 15
					err3 := tx.Save(&user).Error
					if err3 != nil {
						tx.Rollback()
						return err3
					}
				}
			}
		}
		break
	case "3" :
		//healing
		var userother statements.UserOther
		result := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).First(&userother)
		if result.Error != nil {
			return result.Error
		}
		if userother.Lo3 != 1 {
			t := time.Now()
			t_zero := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
			t_to_tomorrow := 24*60*60 - (t.Unix() - t_zero)		
			redis_cli := setting.RedisClient
			healed := !redis_cli.SetNX(fmt.Sprintf("finish_lo3_id:%d", userID), 0, time.Duration(t_to_tomorrow)*time.Second).Val()
			
			if !healed {
				err2 := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).Update("lo3", 1).Error
				if err2 != nil {
					tx.Rollback()
					return err2
				}
				
				var user statements.User
				result := tx.Model(&statements.User{}).Where("id= ?", userID).First(&user)
				if result.Error != nil {
					return result.Error
				}
				if user.Money >= 0 {
					user.Money = user.Money + 20
					err3 := tx.Save(&user).Error
					if err3 != nil {
						tx.Rollback()
						return err3
					}
				}
			}
		}
		break
	case "4" :
		//singHome
		var userother statements.UserOther
		result := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).First(&userother)
		if result.Error != nil {
			return result.Error
		}
		if userother.Lo4 != 1 {
			t := time.Now()
			t_zero := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
			t_to_tomorrow := 24*60*60 - (t.Unix() - t_zero)		
			redis_cli := setting.RedisClient
			sang := !redis_cli.SetNX(fmt.Sprintf("finish_lo4_id:%d", userID), 0, time.Duration(t_to_tomorrow)*time.Second).Val()
			
			if !sang {
				err2 := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).Update("lo4", 1).Error
				if err2 != nil {
					tx.Rollback()
					return err2
				}
				
				var user statements.User
				result := tx.Model(&statements.User{}).Where("id= ?", userID).First(&user)
				if result.Error != nil {
					return result.Error
				}
				if user.Money >= 0 {
					user.Money = user.Money + 20
					err3 := tx.Save(&user).Error
					if err3 != nil {
						tx.Rollback()
						return err3
					}
				}
			}
		}
		break
	case "5" :
		//praise
		var userother statements.UserOther
		result := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).First(&userother)
		if result.Error != nil {
			return result.Error
		}
		if userother.Lo5 != 1 {
			t := time.Now()
			t_zero := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
			t_to_tomorrow := 24*60*60 - (t.Unix() - t_zero)		
			redis_cli := setting.RedisClient
			praised := !redis_cli.SetNX(fmt.Sprintf("finish_lo5_id:%d", userID), 0, time.Duration(t_to_tomorrow)*time.Second).Val()
			
			if !praised {
				err2 := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).Update("lo5", 1).Error
				if err2 != nil {
					tx.Rollback()
					return err2
				}
				
				var user statements.User
				result := tx.Model(&statements.User{}).Where("id= ?", userID).First(&user)
				if result.Error != nil {
					return result.Error
				}
				if user.Money >= 0 {
					user.Money = user.Money + 10
					err3 := tx.Save(&user).Error
					if err3 != nil {
						tx.Rollback()
						return err3
					}
				}
			}
		}
		break
	case "6" :
		//share
		var userother statements.UserOther
		result := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).First(&userother)
		if result.Error != nil {
			return result.Error
		}
		if userother.Lo6 != 1 {
			t := time.Now()
			t_zero := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
			t_to_tomorrow := 24*60*60 - (t.Unix() - t_zero)		
			redis_cli := setting.RedisClient
			shared := !redis_cli.SetNX(fmt.Sprintf("finish_lo6_id:%d", userID), 0, time.Duration(t_to_tomorrow)*time.Second).Val()
			
			if !shared {
				err2 := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).Update("lo6", 1).Error
				if err2 != nil {
					tx.Rollback()
					return err2
				}
				
				var user statements.User
				result := tx.Model(&statements.User{}).Where("id= ?", userID).First(&user)
				if result.Error != nil {
					return result.Error
				}
				if user.Money >= 0 {
					user.Money = user.Money + 10
					err3 := tx.Save(&user).Error
					if err3 != nil {
						tx.Rollback()
						return err3
					}
				}
			}
		}
		break
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
	// log.Println(prop)

	var black GetPrize

	var prize []Prize
	err := tx.Model(&statements.Prize{}).Select("id, name, photo, weight, count").Scan(&prize).Error
	if err != nil {
		log.Println(err)
	}

	var user statements.User
	result0 := tx.Model(&statements.User{}).Where("id= ?", userID).First(&user)
	if result0.Error != nil {
		log.Println(1)
		return black, result0.Error
	}

	// log.Println(user.Money)

	if user.Money >= 100 {
		// t := time.Now()
		// t_zero := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
		// t_to_tomorrow := 24*60*60 - (t.Unix() - t_zero)		
		
		t_to_duration := 1
		redis_cli := setting.RedisClient
		draw_lottery := !redis_cli.SetNX(fmt.Sprintf("draw_lottery_id:%d", userID), 0, time.Duration(t_to_duration)*time.Second).Val()

		if !draw_lottery {
			LOOP:
			for i := 0; i < len(prize); i++ {
				if prop <= prize[i].Weight && prize[i].Count > 0 {
					responseLottery = GetPrize{
						Id:    prize[i].Id,
						Name:  prize[i].Name,
						Photo: prize[i].Photo,
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
	
					//增加匿名次数
					if prize[i].Id == 6 {
						var anonymous statements.UserOther
						result1 := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).First(&anonymous)
						if result1.Error != nil {
							log.Println(2)
							return black, result1.Error
						}
						anonymous.RemainHideName = anonymous.RemainHideName + 1
						err0 := tx.Save(&anonymous).Error
						if err0 != nil {
							tx.Rollback()
						}
					}
	
					//增加点歌次数
					if prize[i].Id == 7 {
						var vod_count statements.UserOther
						result1 := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).First(&vod_count)
						if result1.Error != nil {
							log.Println(3)
							return black, result1.Error
						}
						vod_count.RemainSing = vod_count.RemainSing + 1
						err0 := tx.Save(&vod_count).Error
						if err0 != nil {
							tx.Rollback()
						}
					}
	
					//修改剩余数量
					var count statements.Prize
					result3 := tx.Model(&statements.Prize{}).Where("id = ?", uint(prize[i].Id)).First(&count)
					if result3.Error != nil {
						log.Println(4)
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
						log.Println(5)
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
	
		}
	} else {
		log.Println(result0)
	}
	return responseLottery, tx.Commit().Error
}
