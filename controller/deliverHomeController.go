package controller

import (
	"healing2020/models"

	"github.com/gin-gonic/gin"
)

func AllDeliver(c *gin.Context) {
	Type := c.Query("Type")
	dev, _ := models.DeliverHome(Type)
	c.JSON(200, dev)
}
