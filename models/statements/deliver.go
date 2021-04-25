package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Deliver struct {
	gorm.Model
	UserId    uint   `gorm:"default: 0"`
	Type      int    `gorm:"default: 0"`
	TextField string `gorm:"default: ''"`
	Photo     string `gorm:"default: ''"`
	Record    string `gorm:"default: ''"`
	Praise    int    `gorm:"default: 0"`
}

func DeliverInit() {
	db := setting.MysqlConn()
	if !db.HasTable(&Deliver{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Deliver{}).Error; err != nil {
			panic(err)
		}
		fmt.Println("Table Deliver has been created")
	} else {
		db.AutoMigrate(&Deliver{})
		fmt.Println("Table Deliver has existed")
	}
}
