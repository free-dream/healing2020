package statements

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "healing2020/pkg/setting"
)

type Lottery struct {
    gorm.Model
    PrizeId uint
    UserId uint
    Weight int
}

func LotteryInit() {
    db := setting.MysqlConn()
    if !db.HasTable(&Lottery{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Lottery{}).Error; err != nil {
            panic(err)
        }
        fmt.Println("Table Lottery has been created")
    }else {
        db.AutoMigrate(&Lottery{})
        fmt.Println("Table Lottery has existed")
    }
}

