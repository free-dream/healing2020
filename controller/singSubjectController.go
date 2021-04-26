package controller

import (
	"healing2020/models"

	"github.com/gin-gonic/gin"
)

func SingSubject(c *gin.Context) {
	sub, _ := models.SingSubject()
	c.JSON(200, sub)
}
