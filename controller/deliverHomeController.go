package controller

import (
	"fmt"
	"healing2020/models"
	"healing2020/pkg/e"

	"github.com/gin-gonic/gin"
)

func AllDeliver(c *gin.Context) {
	Type := c.Query("Type")
	user, err := models.DeliverHome(Type)
	if err != nil {
		fmt.Println(err)
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, user)
}
