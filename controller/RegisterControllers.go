package controller

import (
	"healing2020/models"
	"healing2020/models/statements"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"

	"github.com/gin-gonic/gin"
)

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
	userInf := tools.GetUser()
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
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.ERROR_USER_CREATE_FAIL)})
		return
	}
	err = models.CreateBackground(userInf.ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.ERROR_USER_CREATE_FAIL)})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: e.GetMsg(e.SUCCESS)})
}
