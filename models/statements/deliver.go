package statements

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "healing2020/pkg/setting"
)

type Deliver struct {
    gorm.Model
    UserId int
    Type int
    TextField string
    Photo string
    Record string
    Praise int
}

func DeliverInit() {
    db := setting.MysqlConn()
    if !db.HasTable(&Deliver{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Deliver{}).Error; err != nil {
            panic(err)
        }
        fmt.Println("Table Deliver has been created")
    }else {
        db.AutoMigrate(&Deliver{})
        fmt.Println("Table Deliver has existed")
    }
}

