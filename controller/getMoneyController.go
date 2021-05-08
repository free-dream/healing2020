package controller

import (
	"healing2020/models"
	"healing2020/models/statements"
	"healing2020/pkg/e"
	"healing2020/pkg/setting"
	"healing2020/pkg/tools"
	"log"

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

// func UseMoney(c *gin.Context) {
// 	userInf := tools.GetUser(c)

// 	err := models.UseMoney(userInf.ID)
// 	if err != nil {
// 		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
// 		return
// 	}
// 	c.JSON(200, "抽奖成功")
// }

func EarnMoney(c *gin.Context) {
	userInf := tools.GetUser(c)

	err := models.EarnMoney(userInf.ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, "提交任务成功")
}

func GetTask(c *gin.Context) {
	userInf := tools.GetUser(c)

	Task, err := models.GetTask(userInf.ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, Task)
}

func PostQRcode(c *gin.Context) {
	User_id := c.Query("user_id")
	err := models.PostQRcode(User_id)
	if err != nil {
		log.Println(err)
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to add money"})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: "ok"})
}

func LotteryDraw(c *gin.Context) {
	db := setting.MysqlConn()
	userInf := tools.GetUser(c)

	var Background statements.UserOther
	result2 := db.Select("ava_background").Where("user_id = ?", userInf.ID).First(&Background)
	if result2.Error != nil {
		return
	}
	bd := SplitAvaBackgroundtoI(Background.AvaBackground)
	bdstr := Background.AvaBackground

	lot, err := models.LotteryDraw(userInf.ID, bd, bdstr)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, lot)
}
