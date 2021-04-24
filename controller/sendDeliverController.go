package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"

	"github.com/gin-gonic/gin"
)

func PostDeliver(c *gin.Context) {
	UserId := c.Query("userId")
	TextField := c.Query("textField")
	Photo := c.Query("photo")
	Record := c.Query("record")
	err := models.PostDeliver(UserId, TextField, Photo, Record)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to add deliver"})
	}
	c.JSON(200, e.ErrMsgResponse{Message: "发送投递成功！"})
}
