package setting

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"healing2020/pkg/tools"
)

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
	dbName := tools.GetConfig("mysql", "dbName")
	user := tools.GetConfig("mysql", "user")
	password := tools.GetConfig("mysql", "password")
	port := tools.GetConfig("mysql", "port")
	dbInfo := user + ":" + password + "@tcp(" + port + ")/" + dbName + "?charset=utf8&parseTime=True&loc=Local"

	//connect
	db, _ := gorm.Open("mysql", dbInfo)
	db.SingularTable(true)

	return db
}
