package statements

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "healing2020/pkg/setting"
)

type Vod struct {
    gorm.Model
    UserId int
    More string
    Name string
    Singer string
    Style string
    Language string
}

func VodInit() {
    db := setting.MysqlConn()
    if !db.HasTable(&Vod{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Vod{}).Error; err != nil {
            panic(err)
        }
        fmt.Println("Table Vod has been created")
    }else {
        db.AutoMigrate(&Vod{})
        fmt.Println("Table Vod has existed")
    }
}

