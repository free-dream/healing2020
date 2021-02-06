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
    TrueName string
    More string
    Campus string
    Avatar string
    Phone string
    Sex int
    Hoppy string
    Money int
    Setting1 int
    Setting2 int
    Setting3 int
}

func UserInit() {
    db := setting.MysqlConn()
    if !db.HasTable(&User{}) {
        if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&User{}).Error; err != nil {
            panic(err)
        }
        fmt.Println("Table User has been created")
    }else {
        db.AutoMigrate(&User{})
        fmt.Println("Table User has existed")
    }
}
