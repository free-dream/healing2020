package setting

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"healing2020/pkg/tools"
)

var DB *gorm.DB

func init() {
	dbName := tools.GetConfig("mysql", "dbName")
	user := tools.GetConfig("mysql", "user")
	password := tools.GetConfig("mysql", "password")
	port := tools.GetConfig("mysql", "port")
	dbInfo := user + ":" + password + "@tcp(" + port + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s"

	//connect
	var err error
	DB, err = gorm.Open("mysql", dbInfo)
	DB.SingularTable(true)
	DB.DB().SetMaxOpenConns(50)
	DB.DB().SetMaxIdleConns(20)
	if err != nil {
		panic(err)
	}
}

func MysqlConnTest() {
	dbName := tools.GetConfig("mysql", "dbName")
	user := tools.GetConfig("mysql", "user")
	password := tools.GetConfig("mysql", "password")
	port := tools.GetConfig("mysql", "port")
	dbInfo := user + ":" + password + "@tcp(" + port + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	//connect
	db, err := gorm.Open("mysql", dbInfo)
	db.SingularTable(true)
	db.Close()
	if err != nil {
		panic(err)
	}
}

func MysqlConn() *gorm.DB {
	return DB
}
