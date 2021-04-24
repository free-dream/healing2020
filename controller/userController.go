package controller

import (
	"fmt"
	"healing2020/models"
	"healing2020/models/statements"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"

	"github.com/gin-gonic/gin"
)

type PutUserInf struct {
	NickName string `json:"name"`
	More     string `json:"signature"`
	Setting1 int    `json:"setting1"`
	Setting2 int    `json:"setting2"`
	Setting3 int    `json:"setting3"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	TrueName string `json:"truename"`
}

type UserRegister struct {
	NickName string `json:"name"`
	TrueName string `json:"realname"`
	Sex      int    `json:"sex"`
	Phone    string `json:"phone"`
	Campus   string `json:"school"`
}

//@Title Register
//@Description 注册接口
//@Tags user
//@Produce json
//@Param json body UserRegister true "用户注册数据"
//@Router /register [post]
//@Success 200 {object} e.ErrMsgResponse
//@Failure 403 {object} e.ErrMsgResponse
func Register(c *gin.Context) {
	//获取redis用户信息
	userInf := tools.GetUser(c)
	//获取json
	json := UserRegister{}
	c.BindJSON(&json)
	//构建模型
	user := statements.User{
		NickName: json.NickName,
		TrueName: json.TrueName,
		Sex:      json.Sex,
		Phone:    json.Phone,
		Campus:   json.Campus,
	}
	err := models.UpdateUser(user, userInf.ID)
	if err != nil {
		fmt.Println(err)
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.ERROR_USER_CREATE_FAIL)})
		return
	}
	err = models.CreateUserOther(userInf.ID)
	if err != nil {
		fmt.Println(err)
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.ERROR_USER_CREATE_FAIL)})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: e.GetMsg(e.SUCCESS)})
}

//@Title PutUser
//@Description 更新用户信息
//@Tags user
//@Produce json
//@Param json body PutUserInf true "更新的用户信息"
//@Router /user [put]
//@Success 200 {object} e.ErrMsgResponse
//@Failure 403 {object} e.ErrMsgResponse
func PutUser(c *gin.Context) {
	//接受json
	json := PutUserInf{}
	c.BindJSON(&json)
	//获取用户信息
	userInf := tools.GetUser(c)
	//构建模型
	user := statements.User{
		NickName: json.NickName,
		More:     json.More,
		Setting1: json.Setting1,
		Setting2: json.Setting2,
		Setting3: json.Setting3,
		Avatar:   json.Avatar,
		Phone:    json.Phone,
		TrueName: json.TrueName,
	}
	err := models.UpdateUser(user, userInf.ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.ERROR_USER_SAVE_FAIL)})
	} else {
		c.JSON(200, e.ErrMsgResponse{Message: e.GetMsg(e.SUCCESS)})
	}
}
