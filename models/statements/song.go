package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Song struct {
	gorm.Model
	UserId   uint   `gorm:"default: 0"`
	VodId    uint   `gorm:"default: 0"`
	VodSend  uint   `gorm:"default: 0"`
	Name     string `gorm:"default: ''"`
	Praise   int    `gorm:"default: 0"`
	Source   string `gorm:"default: ''"`
	Style    string `gorm:"default: ''"`
	Language string `gorm:"default: ''"`
}

func SongInit() {
	db := setting.MysqlConn()
	if !db.HasTable(&Song{}) {
		if err := db.CreateTable(&Song{}).Error; err != nil {
			panic(err)
		}
		fmt.Println("Table Song has been created")
	} else {
		db.AutoMigrate(&Song{})
		fmt.Println("Table Song has existed")
	}
}
