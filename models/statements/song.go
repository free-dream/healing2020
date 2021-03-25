package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Song struct {
	gorm.Model
	UserId   uint
	VodId    uint
	VodSend  uint
	Name     string
	Praise   int
	Source   string
	Style    string
	Language string
}

func SongInit() {
	db := setting.MysqlConn()
	if !db.HasTable(&Song{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Song{}).Error; err != nil {
			panic(err)
		}
		fmt.Println("Table Song has been created")
	} else {
		db.AutoMigrate(&Song{})
		fmt.Println("Table Song has existed")
	}
}
