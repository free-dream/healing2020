package controller

import (
	"strings"

	"healing2020/models"
	"healing2020/models/statements"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"

	"github.com/gin-gonic/gin"
)

var tagNum int = 7 //标签数
//1.流行 2.古风 3.民谣 4.摇滚 5.抖音热歌 6.acg 7.其它
type Tag struct {
	Tag1 int `json:"tag1"`
	Tag2 int `json:"tag2"`
	Tag3 int `json:"tag3"`
	Tag4 int `json:"tag4"`
	Tag5 int `json:"tag5"`
	Tag6 int `json:"tag6"`
	Tag7 int `json:"tag7"`
}

//拼接要存入数据库的字符串
func hobbyJoin(json Tag) string {
	h := make([]string, tagNum)
	var n int = 0
	if json.Tag1 == 1 {
		h[n] = "1"
		n++
	}
	if json.Tag2 == 1 {
		h[n] = "2"
		n++
	}
	if json.Tag3 == 1 {
		h[n] = "3"
		n++
	}
	if json.Tag4 == 1 {
		h[n] = "4"
		n++
	}
	if json.Tag5 == 1 {
		h[n] = "5"
		n++
	}
	if json.Tag6 == 1 {
		h[n] = "6"
		n++
	}
	if json.Tag7 == 1 {
		h[n] = "7"
		n++
	}
	h = h[:n]
	return strings.Join(h, ",")
}

//分割字符串
func hobbySplit(hobby string) Tag {
	returnH := Tag{
		Tag1: 0,
		Tag2: 0,
		Tag3: 0,
		Tag4: 0,
		Tag5: 0,
		Tag6: 0,
		Tag7: 0,
	}
	//分割字符串
	splitH := strings.Split(hobby, ",")
	len := len(splitH)
	//对结果进行判断
	for i := 0; i < len; i++ {
		switch splitH[i] {
		case "1":
			returnH.Tag1 = 1
		case "2":
			returnH.Tag2 = 1
		case "3":
			returnH.Tag3 = 1
		case "4":
			returnH.Tag4 = 1
		case "5":
			returnH.Tag5 = 1
		case "6":
			returnH.Tag6 = 1
		case "7":
			returnH.Tag7 = 1
		}
	}
	return returnH
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
	hobby := hobbyJoin(json)
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
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
	} else {
		c.JSON(200, hobbySplit(hobby))
	}
}
