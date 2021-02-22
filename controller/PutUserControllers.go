package controller

import (
	"healing2020/models/statements"
	"healing2020/models"
	"healing2020/pkg/tools"
	"healing2020/pkg/e"

	"github.com/gin-gonic/gin"
)

type PutUserInf struct {
	NickName string `json:"name"`
	More string `json:"signature"`
	Setting1 int `json:"setting1"`
	Setting2 int `json:"setting2"`
	Setting3 int `json:"setting3"`
	Avatar string `json:"avatar"`
}

//@Title PutUser
//@Description 更新用户信息
//@Tags PutUserInf
//@Produce json
//@Router /user [put]
//@Success 200 {object} e.ErrMsgResponse
//@Failure 403 {object} e.ErrMsgResponse
func PutUser(c *gin.Context) {
	//接受json
	json := PutUserInf{}
	c.BindJSON(&json) 
	//获取用户信息
	userInf := tools.GetUser()
	//构建模型
	user := statements.User{
		NickName: json.NickName,
		More: json.More,
		Setting1: json.Setting1,
		Setting2: json.Setting2,
		Setting3: json.Setting3,
		Avatar: json.Avatar,
	}
	err := models.UpdateUser(user, userInf.ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.ERROR_USER_SAVE_FAIL)})
	}else{	
		c.JSON(200, e.ErrMsgResponse{Message: e.GetMsg(e.SUCCESS)})
	}
}

