package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"
	"log"

	"github.com/gin-gonic/gin"
)

type PrizeParams struct {
	Name   string `json:"name"`
	Photo  string `json:"photo"`
	Intro  string `json:"intro"`
	Weight int    `json:"weight"`
	Count  int    `json:"count"`
}

func PostPrize(c *gin.Context) {
	var params PrizeParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(400, e.ErrMsgResponse{Message: "Uncomplete params"})
		return
	}
	err := models.GodAddPrize(params.Name, params.Photo, params.Intro, params.Weight, params.Count)
	if err != nil {
		log.Println(err)
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to add prize"})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: "发送奖品成功！"})
}
