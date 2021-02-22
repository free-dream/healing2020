package controller

import (
    "healing2020/models"
	"github.com/gin-gonic/gin"
    "healing2020/pkg/e"
)

func UserRank(c *gin.Context) {

}

func SongsRank(c *gin.Context) {

}

// @Title GetDeliverRank
// @Description 投递页排行榜
// @Tags rank
// @Produce json
// @Router /deliver/rank
// @Success 200 {object} []AllRank
// @Failure 403 {object} e.ErrMsgResponse
func DeliverRank(c *gin.Context) {
    data,err := models.GetDeliverRank()
    if err != "" {
        c.JSON(403,e.ErrMsgResponse{Message:err})
        return
    }
    c.JSON(200,data)
    return
}

type NewInputParams struct {
    UserId uint `json:"UserId" binding:"required"`
    Types int `json:"Types" binding:"required"`
    Textfield string `json:"Textfield" binding:"required"`
    Photo string `json:"Photo" binding:"required"`
    Record string `json:"Record" binding:"required"`
    Praise int `json:"Praise" binding:"required"`
}

func NewDeliverRank(c *gin.Context) {
    var params NewInputParams    
    if err := c.ShouldBind(&params);err != nil {
        c.JSON(400,e.ErrMsgResponse{Message:"Incomplete params"})
        return
    }
     err := models.CreateDeliver(params.UserId,params.Types,params.Textfield,params.Photo,params.Record,params.Praise)
     if err!="" {
         c.JSON(400,e.ErrMsgResponse{Message:"Error print on console"})
     }
     c.JSON(200,e.ErrMsgResponse{Message:"ok"})
}
