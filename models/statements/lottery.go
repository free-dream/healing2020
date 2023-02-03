package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Lottery struct {
	gorm.Model
	PrizeId uint `gorm:"default: 0"`
	UserId  uint `gorm:"default: 0"`
}

func LotteryInit() {
	db := setting.MysqlConn()
	if !db.HasTable(&Lottery{}) {
		if err := db.CreateTable(&Lottery{}).Error; err != nil {
			panic(err)
		}
		fmt.Println("Table Lottery has been created")
	} else {
		db.AutoMigrate(&Lottery{})
		fmt.Println("Table Lottery has existed")
	}
}
