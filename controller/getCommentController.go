package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"

	"github.com/gin-gonic/gin"
)

func GetComment(c *gin.Context) {
	strID := c.Query("id")
	Type := c.Query("Type")

	com, err := models.GetComment(strID, Type)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, com)
	return
}
