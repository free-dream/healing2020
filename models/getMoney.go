package models

import (
	"fmt"
	"healing2020/pkg/setting"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Money struct {
	Money int `json:"money"`
}

func GetMoney(userID uint) ([]Money, error) {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	//获取个人积分信息
	var user []Money
	err := db.Table("User").Select("money").Where("id= ? ", userID).First(&user).Error
	fmt.Println(userID)
	return user, err
}
