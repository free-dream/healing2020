package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Comment struct {
	gorm.Model
	UserId    uint   `gorm:"default: 0"`
	Type      int    `gorm:"default: 0"`
	SongId    uint   `gorm:"default: 0"`
	DeliverId uint   `gorm:"default: 0"`
	Content   string `gorm:"default: ''"`
	Praise	  int    `gorm:"default: 0"`
}

func CommentInit() {
	db := setting.MysqlConn()
	if !db.HasTable(&Comment{}) {
		if err := db.CreateTable(&Comment{}).Error; err != nil {
			panic(err)
		}
		fmt.Println("Table Comment has been created")
	} else {
		db.AutoMigrate(&Comment{})
		fmt.Println("Table Comment has existed")
	}
}
