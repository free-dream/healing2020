package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Prize struct {
	Id     int    `json:"prize_id" `
	Name   string `json:"name"`
	Photo  string `json:"photo"`
	Weight int    `json:"weight"`
}
type LotteryId struct {
	PrizeId uint `json:"prize_id" `
}
type AllLottery struct {
	Id    uint   `json:"prize_id" `
	Name  string `json:"lottery_name"`
	Photo string `json:"lottery_photo"`
}

func AllPrize() ([]Prize, error) {
	//连接mysql
	db := setting.MysqlConn()

	//获取所有奖品信息
	var prizeHome []Prize
	err := db.Table("prize").Select("id, name, photo, weight").Scan(&prizeHome).Error
	return prizeHome, err
}

func MyLottery(userID uint) ([]AllLottery, error) {
	//连接mysql
	db := setting.MysqlConn()

	//获取用户抽奖奖品id
	var UserLotteryId []LotteryId
	err := db.Table("lottery").Select("prize_id").Where("user_id = ?", userID).Scan(&UserLotteryId).Error

	//通过奖品id获取奖品内容
	UserLottery := make([]statements.Prize, len(UserLotteryId))
	for i := 0; i < len(UserLotteryId); i++ {
		err = db.Table("prize").Select("id, name, photo").Where("id = ?", UserLotteryId[i].PrizeId).Scan(&UserLottery[i]).Error
		if err != nil {
			return nil, err
		}
	}

	responseLottery := make([]AllLottery, len(UserLotteryId))
	for i := 0; i < len(UserLotteryId); i++ {
		responseLottery[i] = AllLottery{
			Id:    UserLottery[i].ID,
			Name:  UserLottery[i].Name,
			Photo: UserLottery[i].Photo,
		}
	}

	return responseLottery, err
}
