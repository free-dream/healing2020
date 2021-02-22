package statements

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "healing2020/pkg/setting"
)

type Praise struct {
    gorm.Model
    UserId uint
    Type int    //1:"deliver"  2:"song"  3:"special"
    PraiseId uint
}

func PraiseInit() {
    db := setting.MysqlConn()
    if !db.HasTable(&Praise{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Praise{}).Error; err != nil {
            panic(err)
        }
        fmt.Println("Table Praise has been created")
    }else {
        db.AutoMigrate(&Praise{})
        fmt.Println("Table Praise has existed")
    }
}