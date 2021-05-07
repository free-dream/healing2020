package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"

	"github.com/gin-gonic/gin"
)

type ToSaveBackground struct {
	Background int
}

//@Title ChangeBackground
//@Description 修改个人背景
//@Tags user
//@Produce json
//@Param json body ToSaveBackground true "修改后的个人背景"
//@Router /api/user/background [put]
//@Success 200 {object} e.ErrMsgResponse
//@Failure 403 {object} e.ErrMsgResponse
func ChangeBackground(c *gin.Context) {
	userInf := tools.GetUser(c)

	jsonInf := ToSaveBackground{}
	c.BindJSON(&jsonInf)

	err := models.UpdateUserOtherNow(userInf.ID, jsonInf.Background)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.ERROR_USER_SAVE_FAIL)})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: e.GetMsg(e.SUCCESS)})
}

type RemainNum struct {
	RemainSing int `json:"remainSing"`
	RemainHide int `json:"remainHide"`
}

//@Title GetRemainNum
//@Description 获取登录用户剩余点歌和匿名次数
//@Tags user
//@Produce json
//@Router /api/user/remainNum [get]
//@Success 200 {object} RemainNum
//@Failure 403 {object} e.ErrMsgResponse
func GetRemainNum(c *gin.Context) {
	//获取redis用户信息
	userID := tools.GetUser(c).ID
	userOther, err := models.SelectRemainNum(userID)
	responseInf := RemainNum{
		RemainSing: userOther.RemainSing,
		RemainHide: userOther.RemainHideName,
	}
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
	} else {
		c.JSON(200, responseInf)
	}
}
