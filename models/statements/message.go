package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Message struct {
	gorm.Model
	Send    uint
	Receive uint
	Type    int
	Content string
}

func MessageInit() {
	db := setting.MysqlConn()
	if !db.HasTable(&Message{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Message{}).Error; err != nil {
			panic(err)
		}
		fmt.Println("Table Message has been created")
	} else {
		db.AutoMigrate(&Message{})
		fmt.Println("Table Message has existed")
	}
}
