package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type UserOther struct {
	gorm.Model
	UserId         uint   `gorm:"default: 0"`
	Now            int    `gorm:"default: 1"` //初始默认为B1且只有B1可用
	AvaBackground  string `gorm:"default: '1'"`
	RemainSing     int
	RemainHideName int
}

func UserOtherInit() {
	db := setting.MysqlConn()
	if !db.HasTable(&UserOther{}) {
		if err := db.CreateTable(&UserOther{}).Error; err != nil {
			panic(err)
		}
		fmt.Println("Table UserOther has been created")
	} else {
		db.AutoMigrate(&UserOther{})
		fmt.Println("Table UserOther has existed")
	}
}
