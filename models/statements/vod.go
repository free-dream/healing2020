package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Vod struct {
	gorm.Model
	UserId   uint   `gorm:"default: 0"`
	More     string `gorm:"default: ''"`
	Name     string `gorm:"default: ''"`
	Singer   string `gorm:"default: ''"`
	Style    string `gorm:"default: ''"`
	Language string `gorm:"default: ''"`
	HideName int    `gorm:"default: 0"`
}

func VodInit() {
	db := setting.MysqlConn()
	if !db.HasTable(&Vod{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Vod{}).Error; err != nil {
			panic(err)
		}
		fmt.Println("Table Vod has been created")
	} else {
		db.AutoMigrate(&Vod{})
		fmt.Println("Table Vod has existed")
	}
}
