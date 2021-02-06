package statements

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "healing2020/pkg/setting"
)

type Song struct {
    gorm.Model
    UserId int
    VodId int
    VodSend int
    Name string
    Like int
    Source string
    Style string
    Language string
}

func SongInit() {
    db := setting.MysqlConn()
    if !db.HasTable(&Song{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Song{}).Error; err != nil {
            panic(err)
        }
        fmt.Println("Table Song has been created")
    }else {
        db.AutoMigrate(&Song{})
        fmt.Println("Table Song has existed")
    }
}

