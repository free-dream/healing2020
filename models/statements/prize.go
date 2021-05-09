package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Prize struct {
	gorm.Model
	Name   string `gorm:"default: ''"`
	Intro  string `gorm:"default: ''"`
	Photo  string `gorm:"default: ''"`
	Weight int    `gorm:"default: 0"`
	Count  int 	  `gorm:"default: 0"`
}

func PrizeInit() {
	db := setting.MysqlConn()
	if !db.HasTable(&Prize{}) {
		if err := db.CreateTable(&Prize{}).Error; err != nil {
			panic(err)
		}
		fmt.Println("Table Prize has been created")
	} else {
		db.AutoMigrate(&Prize{})
		fmt.Println("Table Prize has existed")
	}
}
