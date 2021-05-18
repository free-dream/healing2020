package controller

import (
	"healing2020/models"
	"healing2020/pkg/tools"

    "github.com/gin-gonic/gin"
)

func FinishData(c *gin.Context) {
    data := models.GetFinish(tools.GetUser(c).ID)
    c.JSON(200,data)
}
