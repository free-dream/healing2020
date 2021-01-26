package controller

import (
    "healing2020/pkg/tools"
	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
    value := tools.GetConfig("test","appid")
    c.JSON(200,gin.H{"test":value})
}

