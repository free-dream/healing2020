package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"
	"log"

	"github.com/gin-gonic/gin"
)

func AllDeliver(c *gin.Context) {
	Type := c.Query("Type")
	// dev, _ := models.DeliverHome(Type, tools.GetUser(c).ID)
	pageStr := c.Query("pageStr")
	dev, _ := models.DeliverHome(pageStr, Type, tools.GetUser(c).ID)
	c.JSON(200, dev)
}

func SingleDeliver(c *gin.Context) {
	DevId := c.Query("deliver_id")
	sin, err := models.SingleDeliver(DevId, tools.GetUser(c).ID)
	if err != nil {
		log.Println(err)
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to get deliver"})
		return
	}
	c.JSON(200, sin)
}
