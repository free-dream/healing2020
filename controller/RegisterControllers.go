package controller

import (
	"healing2020/models/statements"

	"healing2020/models"
	"healing2020/pkg/tools"

	"github.com/gin-gonic/gin"
)

//@Title Register
//@Description 注册接口
//@Tags register
//@Produce json
//@Router /register [post]
//@Success 200 {string} string "{"message": "xxxxx"}"
//@Failure 403 {string} string "{"error": "false"}"
type UserRegister struct {
	NickName string `json:"name"`
	TrueName string `json:"realname"`
	Sex int `json:"sex"`
	Phone string `json:"phone"`
	Campus string `json:"school"`
}

func Register(c *gin.Context) {
	//获取redis用户信息
	userInf := tools.GetUser() 
	//获取json
	json := UserRegister{}
	c.BindJSON(&json)
	user := statements.User{
		NickName: json.NickName,
		TrueName: json.TrueName,
		Sex: json.Sex,
		Phone: json.Phone,
		Campus: json.Campus,
	}
	err := models.RegisterUpdate(user, userInf.ID)
	if err != nil {
		c.JSON(403, gin.H{"error": err})
	}else{	
		c.JSON(200, gin.H{"message": "注册成功"})
	}
}