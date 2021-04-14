package statements

import (
	"fmt"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Rank struct {
	gorm.Model
	Campus   string `gorm:"default: ''"`
	AllRank  string `gorm:"default: ''"`
	PartRank string `gorm:"default: ''"`
}

func RankInit() {
	db := setting.MysqlConn()
	if !db.HasTable(&Rank{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Rank{}).Error; err != nil {
			panic(err)
		}
		fmt.Println("Table Rank has been created")
	} else {
		db.AutoMigrate(&Rank{})
		fmt.Println("Table Rank has existed")
	}
}
