package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Mailbox struct {
	gorm.Model
	Message string `gorm:"default: ''"`
}

func MailboxInit() {
	db := setting.MysqlConn()
	if !db.HasTable(&Mailbox{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Mailbox{}).Error; err != nil {
			panic(err)
		}
		fmt.Println("Table Mailbox has been created")
	} else {
		db.AutoMigrate(&Mailbox{})
		fmt.Println("Table Mailbox has existed")
	}
}
