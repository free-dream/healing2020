package controller

import (
	"strconv"
	"fmt"

	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/models/statements"
	"healing2020/pkg/tools"

	"github.com/gin-gonic/gin"
)

type PersonalPage struct {
	NickName string `json:"name"`
	Campus string `json:"school"`
	More string  `json:"more"`
	Setting1 int `json:"setting1"`
	Setting2 int `json:"setting2"`
	Setting3 int `json:"setting3"`
	Avatar string `json:"avatar"`
	Background string `json:"background"`
	Vod[] models.RequestSongs `json:"requestSongs"`
	Songs[] models.Songs `json:"Songs"`
	Praise[] models.Admire `json:"admire"` 
}

//综合处理各项数据获取最终返回结果
//type: 1为其它用户个人页接口使用，2为登录用户个人页接口使用
func responsePage(c *gin.Context, user statements.User, userID uint) {
	var err error

	//初始化返回数据
	page := PersonalPage{
		NickName: user.NickName,
		Campus: user.Campus,
		More: user.More,
		Setting1: user.Setting1,
		Setting2: user.Setting2,
		Setting3: user.Setting3,
		Avatar: user.Avatar,
	}

	//补充返回数据
	page.Background, err = models.ResponseBackground(userID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}

	page.Vod, err = models.ResponseVod(userID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}

	page.Songs, err = models.ResponseSongs(userID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}

	page.Praise, err = models.ResponsePraise(userID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, page)
	return
}

//@Title ResponseMyPerponalPage
//@Description 已登录用户的个人页接口
//@Tags my perponalpage
//@Produce json
//@Router /user [get]
//@Success 200 {object} PersonalPage
//@Failure 403 {object} e.ErrMsgResponse
func ResponseMyPerponalPage(c *gin.Context){
	rUser := tools.GetUser()
	user := statements.User(rUser)
	responsePage(c, user, user.ID)
}

//@Title ResponseOthersPerponalPage
//@Description 其它用户的个人页接口
//@Tags others'perponalpage
//@Produce json
//@Router /user/{id} [get]
//@Success 200 {object} PersonalPage
//@Failure 403 {object} e.ErrMsgResponse
func ResponseOthersPerponalPage(c *gin.Context) {
	//获取querystring并转化格式
	id := c.Query("id")
	userIDInt, err := strconv.Atoi(id)
	var userID uint = uint(userIDInt)
	if err != nil{
		fmt.Println("字符串转换成整数失败")
	}
	//查询id对应用户信息
	user, err := models.ResponseUser(userID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	responsePage(c, user, userID)
}
