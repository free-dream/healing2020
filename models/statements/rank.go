package statements

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "healing2020/pkg/setting"
)

type Rank struct {
    gorm.Model
    Campus string
    AllRank string
    PartRank string
}

func RankInit() {
    db := setting.MysqlConn()
    if !db.HasTable(&Rank{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Rank{}).Error; err != nil {
            panic(err)
        }
        fmt.Println("Table Rank has been created")
    }else {
        db.AutoMigrate(&Rank{})
        fmt.Println("Table Rank has existed")
    }
}
