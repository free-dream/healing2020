package statements

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "healing2020/pkg/setting"
)

type Comment struct {
    gorm.Model
    UserId uint
    Type int
    SongId uint
    DeliverId uint
    Content string
}

func CommentInit() {
    db := setting.MysqlConn()
    if !db.HasTable(&Comment{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Comment{}).Error; err != nil {
            panic(err)
        }
        fmt.Println("Table Comment has been created")
    }else {
        db.AutoMigrate(&Comment{})
        fmt.Println("Table Comment has existed")
    }
}

