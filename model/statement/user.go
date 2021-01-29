package statements

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
    gorm.Model
    Openid string
    NickName string
}

func UserInit() {
    if !db.HasTable(&User{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&User{}).Error; err != nil {
            panic(err)
        }
    }else {
        fmt.Println("Table User has existed")
    }
}
