package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	OpenId   string `gorm:"default: ''"`
	NickName string `gorm:"default: ''"`
	TrueName string `gorm:"default: ''"`
	More     string `gorm:"default: ''"`
	Campus   string `gorm:"default: ''"`
	Avatar   string `gorm:"default: ''"`
	Phone    string `gorm:"default: ''"`
	Sex      int    `gorm:"default: 0"`
	Hobby    string `gorm:"default: ''"`
	Money    int    `gorm:"default: 0"`
	Setting1 int    `gorm:"default: 1"`
	Setting2 int    `gorm:"default: 1"`
	Setting3 int    `gorm:"default: 1"`
}

func UserInit() {
	db := setting.MysqlConn()
	if !db.HasTable(&User{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&User{}).Error; err != nil {
			panic(err)
		}
		fmt.Println("Table User has been created")
	} else {
		db.AutoMigrate(&User{})
		fmt.Println("Table User has existed")
	}
}
