package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"

	"github.com/gin-gonic/gin"
)

func PostComment(c *gin.Context) {
	UserId := c.Query("userId")
	Content := c.Query("content")
	id := c.Query("id")
	Type := c.Query("type")
	err := models.PostComment(UserId, id, Type, Content)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to add comment"})
	}
	c.JSON(200, e.ErrMsgResponse{Message: "发送评论成功！"})
}
