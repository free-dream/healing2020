package statements

import (
    "fmt"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "healing2020/pkg/setting"
)

type Background struct {
    ID uint `gorm:"primary_key"`
    UserId uint
    Now int  `gorm:"default: 1"`     //初始默认为B1且只有B1可用
    B1 int   `gorm:"default: 1"`
    B2 int   `gorm:"default: 0"`
    B3 int   `gorm:"default: 0"`
    B4 int  `gorm:"default: 0"`
    B5 int   `gorm:"default: 0"`
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