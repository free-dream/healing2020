package statements

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "healing2020/pkg/setting"
)

type Special struct {
    gorm.Model
    SubjectId uint
    UserId uint
    Name string
    Praise int
    Song string
}

func SpecialInit() {
    db := setting.MysqlConn()
    if !db.HasTable(&Special{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Special{}).Error; err != nil {
            panic(err)
        }
        fmt.Println("Table Special has been created")
    }else {
        db.AutoMigrate(&Special{})
        fmt.Println("Table Special has existed")
    }
}

