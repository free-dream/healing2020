package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"

	"github.com/gin-gonic/gin"
)

func GetMoney(c *gin.Context) {
	userInf := tools.GetUser(c)

	Money, err := models.GetMoney(userInf.ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, Money)
}

func UseMoney(c *gin.Context) {
	userInf := tools.GetUser(c)

	err := models.UseMoney(userInf.ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, "抽奖成功")
}

func EarnMoney(c *gin.Context) {
	userInf := tools.GetUser(c)

	err := models.EarnMoney(userInf.ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, "提交任务成功")
}