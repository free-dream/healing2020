package controller

import (
	"fmt"
	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"
	"strconv"
	"strings"

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
//@Router /api/register [post]
//@Success 200 {object} e.ErrMsgResponse
//@Failure 403 {object} e.ErrMsgResponse
func Register(c *gin.Context) {
	//获取redis用户信息
	userID := tools.GetUser(c).ID
	//获取json
	jsonInf := UserRegister{}
	c.BindJSON(&jsonInf)
	//构建模型
	userMap := map[string]interface{}{
		"NickName": jsonInf.NickName,
		"TrueName": jsonInf.TrueName,
		"Sex":      jsonInf.Sex,
		"Phone":    jsonInf.Phone,
		"Campus":   jsonInf.Campus,
	}
	err := models.UpdateUser(c, userMap, userID)
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
//@Router /api/user [put]
//@Success 200 {object} e.ErrMsgResponse
//@Failure 403 {object} e.ErrMsgResponse
func PutUser(c *gin.Context) {
	//接受json
	jsonInf := PutUserInf{}
	c.BindJSON(&jsonInf)
	//获取用户信息
	userID := tools.GetUser(c).ID
	//构建模型
	userMap := map[string]interface{}{
		"NickName": jsonInf.NickName,
		"More":     jsonInf.More,
		"Setting1": jsonInf.Setting1,
		"Setting2": jsonInf.Setting2,
		"Setting3": jsonInf.Setting3,
		"Phone":    jsonInf.Phone,
		"TrueName": jsonInf.TrueName,
	}
	err := models.UpdateUser(c, userMap, userID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.ERROR_USER_SAVE_FAIL)})
	} else {
		c.JSON(200, e.ErrMsgResponse{Message: e.GetMsg(e.SUCCESS)})
	}
}

type GetUserResp struct {
	ID       int
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

// @Title GetUser
// @Description 获取用户模型，如果path不给id将获取自己的信息
// @Tags user
// @Produce json
// @Router /api/usermodel/{id} [get]
// @Param id path string true "id"
// @Success 200 {object} GetUserResp
// @Failure 401 {object} e.ErrMsgResponse
func GetUser(ctx *gin.Context) {
	if idstr := ctx.Param("id"); idstr == "" {
		user := tools.GetUser(ctx)
		ctx.JSON(200, &user)
		return
	} else {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			ctx.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
			return
		}
		user, err := models.ResponseUser(uint(id))
		user.TrueName = ""
		if err != nil {
			ctx.JSON(404, e.ErrMsgResponse{Message: e.GetMsg(e.NOT_FOUND)})
			return
		}
		ctx.JSON(200, &user)
	}
}

// @Title GetUser
// @Description 获取用户模型，如果path不给id将获取自己的信息
// @Tags user
// @Produce json
// @Router /api/usermodel [get]
// @Success 200 {object} GetUserResp
// @Failure 401 {object} e.ErrMsgResponse
func GetUserExample() {

}

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
	var jsonInf Tag
	c.BindJSON(&jsonInf)
	hobby := hobbyJoin(jsonInf.TagInf)
	//获取redis用户信息
	userID := tools.GetUser(c).ID
	err := models.UpdateUser(c, map[string]interface{}{"Hobby": hobby}, userID)
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
