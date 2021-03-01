package controller

import (
    "healing2020/models"
	"github.com/gin-gonic/gin"
    "healing2020/pkg/e"
    "healing2020/pkg/tools"
)

type PhoneHealing struct {
    Phone string `json:"phone"`
}
// @Title GetUserPhone
// @Description 用户手机
// @Tags heal
// @Produce json
// @Router /user/phone
// @Success 200 {object} PhoneHealing
// @Failure 403 {object} e.ErrMsgResponse
func PhoneHeal(c *gin.Context) {
    data := models.GetPhone() 
    var phoneHealing PhoneHealing
    phoneHealing.Phone = data
    c.JSON(200,phoneHealing)
    return
}

// @Title GetRecord
// @Description 听录音
// @Tags heal
// @Produce json
// @Router /record?id=
// @Success 200 {object} RecordResp
// @Failure 403 {object} e.ErrMsgResponse
func Record(c *gin.Context) {
    id := c.Query("id")
    if !tools.Valid(id,`^[0-9]+$`) {
        c.JSON(403,e.ErrMsgResponse{Message:"Unexpected params"})
        return
    }
    data := models.GetRecord(id)
    if data.Err != nil {
        c.JSON(403,e.ErrMsgResponse{Message:"Fail to get record"})
        return
    }
    c.JSON(200,data)
    return
}

// @Title AddPraise
// @Description 点赞
// @Tags heal
// @Produce json
// @Router /like?id=&type=
// @Success 200 {object} e.ErrMsgResponse
// @Failure 403 {object} e.ErrMsgResponse
func Praise(c *gin.Context) {
    id := c.Query("id")
    types := c.Query("type")
    if !tools.Valid(id,`^[0-9]+$`) || !tools.Valid(types,`^[12]$`) {
        c.JSON(403,e.ErrMsgResponse{Message:"Unexpected Params"})
        return
    }
    err := models.AddPraise(id,types)
    if err != nil {
        c.JSON(403,e.ErrMsgResponse{Message:"Fail to add praise"})
        return 
    }
    c.JSON(200,e.ErrMsgResponse{Message:"ok"})
    return
}

type RecordParams struct {
    Id string `json:"id" binding:"required"`
    Name string `json:"name" binding:"required"`
    Url string `json:"url" binding:"required"`
}
// @Title AddRecord
// @Description 录音治愈发布
// @Tags heal
// @Produce json
// @Router /record
// @Success 200 {object} e.ErrMsgResponse
// @Failure 403 {object} e.ErrMsgResponse
func RecordHeal(c *gin.Context) {
    var params RecordParams
    if err:=c.ShouldBind(&params);err!=nil {
        c.JSON(400,e.ErrMsgResponse{Message:"Uncomplete params"})
        return
    }
    err := models.CreateRecord(params.Id,params.Url)
    if err != nil {
        c.JSON(403,e.ErrMsgResponse{Message:"Fail to add praise"})
    }
    c.JSON(200,e.ErrMsgResponse{Message:"ok"})
} 

type VodParams struct {
    Songs string `json:"songs" binding:"required"`
    Singer string `json:"singer" binding:"required"`
    More string `json:"more" binding:"required"`
    Style string `json:"style" binding:"required"`
    Language string `json:"language" binding:"required"`
}
// @Title AddVod
// @Description 点歌
// @Tags heal
// @Produce json
// @Router /vod
// @Success 200 {object} e.ErrMsgResponse
// @Failure 403 {object} e.ErrMsgResponse
func VodPost(c *gin.Context) {
    var params VodParams
    if err:=c.ShouldBind(&params);err!=nil {
        c.JSON(400,e.ErrMsgResponse{Message:"Uncomplete params"})
        return
    }
    err := models.CreateVod(params.Singer,params.Style,params.Language,params.Songs,params.More)
    if err != nil {
        c.JSON(403,e.ErrMsgResponse{Message:"Fail to add praise"})
    }
    c.JSON(200,e.ErrMsgResponse{Message:"ok"})
}
