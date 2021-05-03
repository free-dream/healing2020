package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type UserOther struct {
	gorm.Model
	UserId         uint `gorm:"default: 0"`
	Now            int  `gorm:"default: 1"` //初始默认为B1且只有B1可用
	B1             int  `gorm:"default: 1"`
	B2             int  `gorm:"default: 0"`
	B3             int  `gorm:"default: 0"`
	B4             int  `gorm:"default: 0"`
	B5             int  `gorm:"default: 0"`
	RemainSing     int  `gorm:"default: 0"`
	RemainHideName int  `gorm:"default: 0"`
	Lo1            int  `gorm:"default: 0"`
	Lo2            int  `gorm:"default: 0"`
	Lo3            int  `gorm:"default: 0"`
	Lo4            int  `gorm:"default: 0"`
	Lo5            int  `gorm:"default: 0"`
	Lo6            int  `gorm:"default: 0"`
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
