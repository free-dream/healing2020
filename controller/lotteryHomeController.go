package controller

import (
	"fmt"
	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"

	"github.com/gin-gonic/gin"
)

func ALLPrize(c *gin.Context) {
	prize, err := models.AllPrize()
	if err != nil {
		fmt.Println(err)
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, prize)
}

func UserLottery(c *gin.Context) {
	userInf := tools.GetUser(c)

	lottery, err := models.MyLottery(userInf.ID)
	if err != nil {
		fmt.Println(err)
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, lottery)
}
