package controller

import (
	"strings"

	"healing2020/models"
	"healing2020/models/statements"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"

	"github.com/gin-gonic/gin"
)

type Tag struct {
	TagInf []string
}

func hobbyJoin(tag []string) string {
	var hobby string
	for key, value := range tag {
		if key == 0 {
			hobby = value
		} else {
			hobby = hobby + "," + value
		}
	}
	return hobby
}

//@Title NewHobby
//@Description 爱好选择接口
//@Tags hobby
//@Produce json
//@Param json body Tag true "用户爱好标签"
//@Router /api/user/hobby [post]
//@Success 200 {object} e.ErrMsgResponse
//@Failure 403 {object} e.ErrMsgResponse
func NewHobby(c *gin.Context) {
	//获取json
	var json Tag
	c.BindJSON(&json)
	hobby := hobbyJoin(json.TagInf)
	//获取redis用户信息
	userInf := tools.GetUser(c)
	err := models.UpdateUser(statements.User{Hobby: hobby}, userInf.ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.ERROR_USER_SAVE_FAIL)})
	} else {
		c.JSON(200, e.ErrMsgResponse{Message: e.GetMsg(e.SUCCESS)})
	}
}

//@Title GetHobby
//@Description 获取用户爱好
//@Tags hobby
//@Produce json
//@Router /api/user/hobby [get]
//@Success 200 {object} Tag
//@Failure 403 {object} e.ErrMsgResponse
func GetHobby(c *gin.Context) {
	//获取redis用户信息
	user := tools.GetUser(c)
	hobby, err := models.SelectUserHobby(user.ID)
	t := Tag{TagInf: strings.Split(hobby, ",")}
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
	} else {
		c.JSON(200, t)
	}
}
