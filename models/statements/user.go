package statements

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "healing2020/pkg/setting"
)

type User struct {
    gorm.Model
    OpenId string
    NickName string
}

func UserInit() {
    db := setting.MysqlConn()
    if !db.HasTable(&User{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&User{}).Error; err != nil {
            panic(err)
        }
    }else {
        fmt.Println("Table User has existed")
    }
}
