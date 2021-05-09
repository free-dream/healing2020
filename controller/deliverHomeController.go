package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"
	"log"

	"github.com/gin-gonic/gin"
)

// type DeliverId struct {
// 	DevId string `json:"deliver_id"`
// }

func AllDeliver(c *gin.Context) {
	Type := c.Query("Type")
	pageStr := c.Query("pageStr")
	dev, _ := models.DeliverHome(pageStr, Type)
	c.JSON(200, dev)
}

func SingleDeliver(c *gin.Context) {
	DevId := c.Query("deliver_id")
	sin, err := models.SingleDeliver(DevId)
	if err != nil {
		log.Println(err)
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to get deliver"})
		return
	}
	c.JSON(200, sin)
}
