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
func GetTask(userID uint) ([]interface{}, error) {
	redis_cli := setting.RedisClient
	taskGet, err := redis_cli.MGet(fmt.Sprintf("%d:%d", userID, 1), fmt.Sprintf("%d:%d", userID, 2), fmt.Sprintf("%d:%d", userID, 3), fmt.Sprintf("%d:%d", userID, 4), fmt.Sprintf("%d:%d", userID, 5), fmt.Sprintf("%d:%d", userID, 6)).Result()
	return taskGet, err
}

//完成每日任务
func FinishTask(task string, userID uint) error {
	//连接mysql
	db := setting.MysqlConn()
	tx := db.Begin()

	var user statements.User
	result2 := tx.Model(&statements.User{}).Where("id= ?", userID).First(&user)
	if result2.Error != nil {
		log.Println(1)
		return result2.Error
	}

	earnMoney := map[string]int{"1": 80, "2": 60, "3": 60, "4": 60, "5": 20, "6": 20}

	t := time.Now()
	t_zero := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	t_to_tomorrow := 24*60*60 - (t.Unix() - t_zero)
	redis_cli := setting.RedisClient

	e := redis_cli.SIsMember(fmt.Sprintf("%d:%s", userID, task), true).Err()
	if e == nil {
		finished := redis_cli.Set(fmt.Sprintf("%d:%s", userID, task), true, time.Duration(t_to_tomorrow)*time.Second).Err()
		if finished == nil {
			if user.Money >= 0 {
				user.Money = user.Money + earnMoney[task]
				err3 := tx.Save(&user).Error
				if err3 != nil {
					log.Println(2)
					tx.Rollback()
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

	var count int

	if user.Money >= 100 {
		t_to_duration := 1
		redis_cli := setting.RedisClient
		draw_lottery := !redis_cli.SetNX(fmt.Sprintf("draw_lottery_id:%d", userID), 0, time.Duration(t_to_duration)*time.Second).Val()

		if !draw_lottery {
			CountFirst := tx.Model(&statements.Lottery{}).Where("user_id = ?", userID).Count(&count).Error
			if CountFirst != nil {
				tx.Rollback()
			}
			//首次抽奖得个人背景
			if count == 0 {
				bdstr = "1,2"
				result2 := tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).Update("ava_background", bdstr).Error
				if result2 != nil {
					tx.Rollback()
				}
				responseLottery = GetPrize{
					Id:    5,
					Name:  "个人页背景",
					Photo: "static/prize5.jpeg",
				}
				//存入我的奖品
				var lot statements.Lottery
				lot.PrizeId = 5
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
			} else {
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
								bdstr = bdstr + "," + strconv.Itoa(len(bd)+1)
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
		}
	} else {
		log.Println(result0)
	}
	return responseLottery, tx.Commit().Error
}
