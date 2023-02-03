package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Praise struct {
	gorm.Model
	UserId   uint `gorm:"default: 0"`
	Type     int  `gorm:"default: 0"` //1:"deliver"  2:"song"  3:"special"
	PraiseId uint `gorm:"default: 0"`
    IsCancel int  `gorm:"default: 0"` //1表示取消点赞，该项作废
}

func PraiseInit() {
	db := setting.MysqlConn()
	if !db.HasTable(&Praise{}) {
		if err := db.CreateTable(&Praise{}).Error; err != nil {
			panic(err)
		}
		fmt.Println("Table Praise has been created")
	} else {
		db.AutoMigrate(&Praise{})
		fmt.Println("Table Praise has existed")
	}
}
