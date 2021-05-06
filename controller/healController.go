package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"

	"github.com/gin-gonic/gin"
)

type PhoneHealing struct {
	Phone string `json:"phone"`
}

// @Title GetUserPhone
// @Description 用户手机
// @Tags heal
// @Produce json
// @Router /api/user/phone [get]
// @Success 200 {object} PhoneHealing
// @Failure 403 {object} e.ErrMsgResponse
func PhoneHeal(c *gin.Context) {
	data := models.GetPhone(tools.GetUser(c))
	var phoneHealing PhoneHealing
	phoneHealing.Phone = data
	c.JSON(200, phoneHealing)
	return
}

type RealResp struct {
	Source string `json:"url"`
}

// @Title GetRecord
// @Description 听录音
// @Tags heal
// @Produce json
// @Router /api/record [get]
// @Param id query string true "record id"
// @Success 200 {object} RealResp
// @Failure 403 {object} e.ErrMsgResponse
func Record(c *gin.Context) {
	id := c.Query("id")
	if !tools.Valid(id, `^[0-9]+$`) {
		c.JSON(403, e.ErrMsgResponse{Message: "Unexpected params"})
		return
	}
	data := models.GetRecord(id)
	if data.Err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to get record"})
		return
	}
	// var realResp RealResp
	// realResp.Source = data.Source
	c.JSON(200, data)
	return
}

// @Title AddPraise
// @Description 点赞
// @Tags heal
// @Produce json
// @Router /api/like [get]
// @Param id query string true "type id"
// @Param type query string true "1 song; 2 deliver"
// @Success 200 {object} e.ErrMsgResponse
// @Failure 403 {object} e.ErrMsgResponse
func Praise(c *gin.Context) {
	id := c.Query("id")
	types := c.Query("type")
	if !tools.Valid(id, `^[0-9]+$`) || !tools.Valid(types, `^[12]$`) {
		c.JSON(403, e.ErrMsgResponse{Message: "Unexpected Params"})
		return
	}
	err := models.AddPraise(id, types)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to add praise"})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: "ok"})
	return
}

type RecordParams struct {
	Id       string   `form:"id" binding:"required"`
	Name     string   `form:"name" binding:"required"`
	ServerID []string `form:"server_id" binding:"required"`
}

// @Title AddRecord
// @Description 录音治愈发布
// @Tags heal
// @Produce json
// @Router /api/record [post]
// @Param id formData string true "点歌单id"
// @Param name formData string false "user name"
// @Param server_id formData []string true "server_id"
// @Success 200 {object} e.ErrMsgResponse
// @Failure 403 {object} e.ErrMsgResponse
func RecordHeal(c *gin.Context) {
	var params RecordParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(400, e.ErrMsgResponse{Message: "Uncomplete params"})
		return
	}
	url, err := convertMediaIdArrToQiniuUrl(params.ServerID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: err.Error()})
		return
	}
	err = models.CreateRecord(params.Id, url, tools.GetUser(c).ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: err.Error()})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: "ok"})
}

type VodParams struct {
	Songs    string `json:"songs" binding:"required"`
	Singer   string `json:"singer"`
	More     string `json:"more"`
	Style    string `json:"style" binding:"required"`
	Language string `json:"language" binding:"required"`
}

// @Title AddVod
// @Description 点歌
// @Tags heal
// @Produce json
// @Router /api/vod [post]
// @Param songs body VodParams true "song's name"
// @Param singer body VodParams false "singer"
// @Param more body VodParams false "备注"
// @Param style body VodParams true "style"
// @Param language body VodParams true "language"
// @Success 200 {object} e.ErrMsgResponse
// @Failure 403 {object} e.ErrMsgResponse
func VodPost(c *gin.Context) {
	var params VodParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(400, e.ErrMsgResponse{Message: "Uncomplete params"})
		return
	}
	err := models.CreateVod(tools.GetUser(c).ID, params.Singer, params.Style, params.Language, params.Songs, params.More)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to add praise"})
	}
	c.JSON(200, e.ErrMsgResponse{Message: "ok"})
}
