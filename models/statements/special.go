package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Special struct {
	gorm.Model
	SubjectId uint   `gorm:"default: 0"`
	UserId    uint   `gorm:"default: 0"`
	Name      string `gorm:"default: ''"`
	Praise    int    `gorm:"default: 0"`
	Song      string `gorm:"default: ''"`
}

func SpecialInit() {
	db := setting.MysqlConn()
	if !db.HasTable(&Special{}) {
		if err := db.CreateTable(&Special{}).Error; err != nil {
			panic(err)
		}
		fmt.Println("Table Special has been created")
	} else {
		db.AutoMigrate(&Special{})
		fmt.Println("Table Special has existed")
	}
}
