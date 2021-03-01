package controller

import (
    "healing2020/models"
	"github.com/gin-gonic/gin"
    "healing2020/pkg/e"
    "healing2020/pkg/tools"
)

// @Title GetUserRank
// @Description 用户积分排行榜
// @Tags rank
// @Produce json
// @Router /user/rank [get]
// @Success 200 {object} []models.AllRank
// @Failure 403 {object} e.ErrMsgResponse
func AllUserRank(c *gin.Context) {
    data,err := models.GetAllUserRank()
    if err != "" {
        c.JSON(403,e.ErrMsgResponse{Message:err})
        return
    }
    c.JSON(200,data)
    return
}


// @Title GetUserRank
// @Description 用户排名
// @Tags rank
// @Produce json
// @Router /user/rank [get]
// @Params id query string
// @Success 200 {object} models.UserRank
// @Failure 403 {object} e.ErrMsgResponse
func UserRank(c *gin.Context) {
    if c.Query("id") == "" {
        AllUserRank(c)
        return
    } 
    id := c.Query("id")
    if tools.Valid(id,`^[0-9]+$`)==false {
        c.JSON(403,e.ErrMsgResponse{Message:"error param"})
        return
    }
    data,err := models.GetUserRank(id)
    if err != nil {
        c.JSON(403,e.ErrMsgResponse{Message:"can not get rank"})
        return
    }
    c.JSON(200,data)
    return
}

// @Title GetSongRank
// @Description 每日歌曲排行榜
// @Tags rank
// @Produce json
// @Router /songs/rank [get]
// @Success 200 {object} []models.AllRank
// @Failure 403 {object} e.ErrMsgResponse
func SongRank(c *gin.Context) {
    data,err := models.GetSongRank()
    if err != "" {
        c.JSON(403,e.ErrMsgResponse{Message:err})
        return
    }
    c.JSON(200,data)
    return
}

// @Title GetDeliverRank
// @Description 投递页排行榜
// @Tags rank
// @Produce json
// @Router /deliver/rank [get]
// @Success 200 {object} []models.AllRank
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
