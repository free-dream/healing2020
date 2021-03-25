package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type UserOther struct {
	ID             uint `gorm:"primary_key"`
	UserId         uint
	Now            int `gorm:"default: 1"` //初始默认为B1且只有B1可用
	B1             int `gorm:"default: 1"`
	B2             int `gorm:"default: 0"`
	B3             int `gorm:"default: 0"`
	B4             int `gorm:"default: 0"`
	B5             int `gorm:"default: 0"`
	RemainSing     int `gorm:"default: 0"`
	RemainHideName int `gorm:"default: 0"`
}

func UserOtherInit() {
	db := setting.MysqlConn()
	if !db.HasTable(&UserOther{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&UserOther{}).Error; err != nil {
			panic(err)
		}
		fmt.Println("Table UserOther has been created")
	} else {
		db.AutoMigrate(&UserOther{})
		fmt.Println("Table UserOther has existed")
	}
}
