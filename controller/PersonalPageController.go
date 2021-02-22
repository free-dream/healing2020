package controller

import (
	"healing2020/models"
	"healing2020/pkg/tools"
	"healing2020/pkg/e"

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

//@Title ResponsePerponalPage
//@Description 个人页接口
//@Tags perponalpage
//@Produce json
//@Router /user [get]
//@Success 200 {object} PersonalPag
//@Failure 403 {object} e.ErrMsgResponse
func ResponsePerponalPage(c *gin.Context){
	var err error

	//获取用户信息
	user := tools.GetUser()
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
	page.Background, err = models.ResponseBackground(user.ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}

	page.Vod, err = models.ResponseVod(user.ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}

	page.Songs, err = models.ResponseSongs(user.ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}

	page.Praise, err = models.ResponsePraise(user.ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	
	c.JSON(200, page)
}