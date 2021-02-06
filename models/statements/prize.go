package statements

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "healing2020/pkg/setting"
)

type Prize struct {
    gorm.Model
    Name string 
    Intro string
    Photo string
    Weight int
}

func PrizeInit() {
    db := setting.MysqlConn()
    if !db.HasTable(&Prize{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Prize{}).Error; err != nil {
            panic(err)
        }
        fmt.Println("Table Prize has been created")
    }else {
        db.AutoMigrate(&Prize{})
        fmt.Println("Table Prize has existed")
    }
}

