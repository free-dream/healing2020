package tools

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type RedisUser struct {
	gorm.Model
	OpenId   string
	NickName string
	TrueName string
	More     string
	Campus   string
	Avatar   string
	Phone    string
	Sex      int
	Hobby    string
	Money    int
	Setting1 int
	Setting2 int
	Setting3 int
}

func GetUser(c *gin.Context) RedisUser {
	session := sessions.Default(c)
	data := session.Get("user")
	return data.(RedisUser)
}
