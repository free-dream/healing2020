package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"

	"github.com/gin-gonic/gin"
)

type DeliverParams struct {
	UserId    string `json:"userId" binding:"required"`
	TextField   string `json:"textField"`
	Photo     string `json:"photo"`
	Record    string `json:"record"`
}

func PostDeliver(c *gin.Context) {
	var params DeliverParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(400, e.ErrMsgResponse{Message: "Uncomplete params"})
		return
	}
	err := models.PostDeliver(params.UserId, params.TextField, params.Photo, params.Record)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to add deliver"})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: "发送投递成功！"})
}
