package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"

  _ "github.com/jinzhu/gorm/dialects/mysql"
)

//SELECT hobby FROM user where id = userID
func HobbySelect(userID uint) (string,error) {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()
	var userHobby statements.User
	err := db.Select("hobby").Where("id=?", userID).Find(&userHobby).Error
	return userHobby.Hobby, err
}