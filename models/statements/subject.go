package statements

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "healing2020/pkg/setting"
)

type Subject struct {
    gorm.Model
    Name string
    Intro string
}

func SubjectInit() {
    db := setting.MysqlConn()
    if !db.HasTable(&Subject{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Subject{}).Error; err != nil {
            panic(err)
        }
        fmt.Println("Table Subject has been created")
    }else {
        db.AutoMigrate(&Subject{})
        fmt.Println("Table Subject has existed")
    }
}

