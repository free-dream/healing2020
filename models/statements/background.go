package statements

import (
    "fmt"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "healing2020/pkg/setting"
)

type Background struct {
    ID uint `gorm:"primary_key"`
    B1 string
    B2 string
    B3 string
    B4 string
    B5 string
}

func BackgroundInit() {
    db := setting.MysqlConn()
    if !db.HasTable(&Background{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Background{}).Error; err != nil {
            panic(err)
        }
        fmt.Println("Table Background has been created")
    }else {
        db.AutoMigrate(&Background{})
        fmt.Println("Table Background has existed")
    }
}