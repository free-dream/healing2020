package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"

	"github.com/gin-gonic/gin"
)

type CommentParams struct {
	UserId  string `json:"userId" binding:"required"`
	Content string `json:"content" binding:"required"`
	Id      string `json:"id" binding:"required"`
	Type    string `json:"type" binding:"required"`
}

func PostComment(c *gin.Context) {
	var params CommentParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(400, e.ErrMsgResponse{Message: "Uncomplete params"})
		return
	}
	err := models.PostComment(params.UserId, params.Id, params.Type, params.Content)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to add comment"})
	}
	c.JSON(200, e.ErrMsgResponse{Message: "发送评论成功！"})
}
