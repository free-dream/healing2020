package models

import (
	"encoding/json"
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"healing2020/pkg/tools"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"errors"
	//"encoding/json"
	//"time"
)

func UpdateOrCreate(openId string, nickName string, sex int, avatar string) {
	db := setting.MysqlConn()
	db.Transaction(func(tx *gorm.DB) error {
		var user statements.User
		result := tx.Model(&statements.User{}).Where("open_id=?", openId).First(&user)
		user.NickName = nickName
		user.OpenId = openId
		user.Avatar = avatar
		user.Sex = sex
		var result2 *gorm.DB
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			var userOther statements.UserOther
			userOther.UserId = user.ID
			tx.Model(&statements.UserOther{}).Create(&userOther)
			result2 = tx.Model(&statements.User{}).Create(&user)
		} else {
			result2 = tx.Model(&statements.User{}).Where("open_id=?", openId).Update(&user)
		}
		// client := setting.RedisConn()
		// defer client.Close()
		// dataByte,_ := json.Marshal(user)
		// data := string(dataByte)
		// keyname := "healing2020:token:"+token
		// client.Set(keyname,data,time.Minute*30)

		return result2.Error
	})
}

//SELECT hobby FROM user where id = userID
func SelectUserHobby(userID uint) (string, error) {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	var userHobby statements.User
	err := db.Select("hobby").Where("id=?", userID).First(&userHobby).Error
	return userHobby.Hobby, err
}

//更新用户数据后要引用一次该函数
func updateSession(c *gin.Context, db *gorm.DB) {
	var redisUser tools.RedisUser
	var user statements.User

	userID := tools.GetUser(c).ID
	db.Where("id=?", userID).First(&user)

	tmp, _ := json.Marshal(user)
	json.Unmarshal(tmp, &redisUser)

	session := sessions.Default(c)
	session.Set("user", redisUser)
	session.Save()
}

//更新user表
func UpdateUser(c *gin.Context, user statements.User, userID uint) error {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	userMap := map[string]interface{}{
		"nick_name": user.NickName,
		"more":      user.More,
		"setting1":  user.Setting1,
		"setting2":  user.Setting2,
		"setting3":  user.Setting3,
		"phone":     user.Phone,
		"true_name": user.TrueName,
	}

	//开启事务
	tx := db.Begin()
	err := tx.Model(&statements.User{}).Where("id=?", userID).Update(userMap).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	com := tx.Commit()
	updateSession(c, db)
	return com.Error
}

//更新user表
func RegisterUpdate(c *gin.Context, user statements.User, userID uint) error {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	//开启事务
	tx := db.Begin()
	err := tx.Model(&statements.User{}).Where("id=?", userID).Update(user).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	com := tx.Commit()
	updateSession(c, db)
	return com.Error
}

//获取用户数
func GetUserNum() (int, error) {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	var count int
	err := db.Table("user").Count(&count).Error
	return count, err
}
