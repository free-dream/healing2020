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
	"strconv"
	//"encoding/json"
	//"time"
)

//run in personalPagecontroller.go
//run in userController.go
//select并根据id返回用户信息
func ResponseUser(userID uint) (statements.User, error) {
	//连接mysql
	db := setting.MysqlConn()

	var user statements.User
	err := db.Where("id=?", userID).First(&user).Error
	return user, err
}

func UpdateOrCreate(openId string, nickName string, sex int, avatar string) error {
	db := setting.MysqlConn()
	err := db.Transaction(func(tx *gorm.DB) error {
		var user statements.User
		result := tx.Model(&statements.User{}).Where("open_id=?", openId).First(&user)
		user.NickName = nickName
		user.OpenId = openId
		user.Avatar = avatar
		user.Sex = sex
		var result2 *gorm.DB
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			result2 = tx.Model(&statements.User{}).Create(&user)
			var userOther statements.UserOther
			userOther.UserId = user.ID
			userOther.RemainSing = 8     //点歌次数默认值
			userOther.RemainHideName = 0 //匿名次数默认值
			tx.Model(&statements.UserOther{}).Create(&userOther)
		} else {
			result2 = tx.Model(&statements.User{}).Where("open_id=?", openId).Update(&user)
		}
		// client := setting.RedisConn()
		// dataByte,_ := json.Marshal(user)
		// data := string(dataByte)
		// keyname := "healing2020:token:"+token
		// client.Set(keyname,data,time.Minute*30)

		return result2.Error
	})
	if err != nil {
		return err
	}
	return nil
}

func UpdatePraiseSign() {
	db := setting.MysqlConn()
	client := setting.RedisConn()

	count := 0
	db.Model(&statements.User{}).Count(&count)

	for i := 0; i < count; i++ {
		keyname := "healing2020:PraiseSign" + strconv.Itoa(count)
		client.Del(keyname)
	}
}

//SELECT hobby FROM user where id = userID
func SelectUserHobby(userID uint) (string, error) {
	//连接mysql
	db := setting.MysqlConn()

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

//踩坑：不要用0值作为状态值, gorm模型不更新0和""，要用map
//更新user表
func UpdateUser(c *gin.Context, userMap map[string]interface{}, userID uint) error {
	//连接mysql
	db := setting.MysqlConn()

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

//获取用户数
func GetUserNum() (int, error) {
	//连接mysql
	db := setting.MysqlConn()

	var count int
	err := db.Table("user").Count(&count).Error
	return count, err
}
