package statements

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "healing2020/pkg/setting"
)

type Message struct {
    gorm.Model
    Send int
    Receive int
    Type int
    Content string
    Url string
}

func MessageInit() {
    db := setting.MysqlConn()
    if !db.HasTable(&Message{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Message{}).Error; err != nil {
            panic(err)
        }
        fmt.Println("Table Message has been created")
    }else {
        db.AutoMigrate(&Message{})
        fmt.Println("Table Message has existed")
    }
}

