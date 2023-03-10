package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"
	"strconv"
	"time"

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
	data := models.GetRecord(id, tools.GetUser(c).ID)
	if data.Err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to get record"})
		return
	}
	// var realResp RealResp
	// realResp.Source = data.Source
	c.JSON(200, data)
	return
}

// @Title CancelPraise
// @Description 取消点赞
// @Tags heal
// @Produce json
// @Router /api/unlike [get]
// @Param id query string true "type id"
// @Param type query string true "1 deliver; 2 song;3 singHome"
// @Success 200 {object} e.ErrMsgResponse
// @Failure 403 {object} e.ErrMsgResponse
func NoPraise(c *gin.Context) {
	id := c.Query("id")
	types := c.Query("type")
	if !tools.Valid(id, `^[0-9]+$`) || !tools.Valid(types, `^[1234]$`) {
		c.JSON(403, e.ErrMsgResponse{Message: "Unexpected Params"})
		return
	}
	err := models.CancelPraise(tools.GetUser(c).ID, id, types)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to add praise"})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: "ok"})
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
	if !tools.Valid(id, `^[0-9]+$`) || !tools.Valid(types, `^[1234]$`) {
		c.JSON(403, e.ErrMsgResponse{Message: "Unexpected Params"})
		return
	}
	err, data := models.AddPraise(tools.GetUser(c).ID, id, types)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to add praise"})
		return
	}
	SendPraiseMsg(data.MyID, data.TargetID, tools.GetUser(c).NickName, data.Type, data.Msg)
	c.JSON(200, e.ErrMsgResponse{Message: "ok"})
	return
}

func SendPraiseMsg(myID uint, targetID uint, myName string, types string, mainMsg string) {
	if types == "2" {
		types = "[治愈]:"
	} else if types == "1" {
		types = "[投递]:"
	}
	content := myName + "点赞了您的" + types + mainMsg
	msg := Message{
		Type:           2,
		Time:           "",
		FromUserID:     0,
		ToUserID:       targetID,
		Content:        content,
		URL:            "",
		IsToFromUserID: 0,
	}
	msgID := tools.Md5String(strconv.Itoa(int(myID)) + strconv.Itoa(int(targetID)) + time.Now().Format("2006-01-02 15:04:05"))
	msg.ID = msgID
	MysqlCreate <- &msg
	createUserMsgChan(targetID) //in ws.go
	MessageQueue[int(targetID)] <- &msg
}

type RecordParams struct {
	Id       string   `json:"id" binding:"required"`
	Name     string   `json:"name"`
	ServerID []string `json:"server_id" binding:"required"`
	IsHide   int      `json:"isHide"`
}

// @Title AddRecord
// @Description 录音治愈发布
// @Tags heal
// @Produce json
// @Router /api/record [post]
// @Param id body string true "点歌单id"
// @Param name body string false "user name"
// @Param server_id body []string true "server_id"
// @Param isHide body int true "1 自己可见,0 所有人可见"
// @Success 200 {object} e.ErrMsgResponse
// @Failure 403 {object} e.ErrMsgResponse
func RecordHeal(c *gin.Context) {
	var params RecordParams
	userID := tools.GetUser(c).ID
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(400, e.ErrMsgResponse{Message: err.Error()})
		return
	}
	if !tools.Valid(strconv.Itoa(params.IsHide), `^[01]$`) {
		c.JSON(400, e.ErrMsgResponse{Message: "Unexpected input"})
		return
	}
	url, err := convertMediaIdArrToQiniuUrl(params.ServerID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: err.Error()})
		return
	}
	songName, err := models.CreateRecord(params.Id, url, userID, params.IsHide)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: err.Error()})
		return
	}

	//send ws song_record message(in ws.go)
	msg := Message{
		Type:       1,
		Time:       time.Now().Format("2006-01-02 15:04:05"),
		FromUserID: userID,
		URL:        url,
	}
	intId, _ := strconv.Atoi(params.Id)
	vodId := uint(intId)
	toUserID, err := models.SelectUserIDByVodID(vodId)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: err.Error()})
		return
	}
	msg.ToUserID = toUserID
	msg.Content = songName
	md5ID := tools.Md5String(strconv.Itoa(int(userID)) + strconv.Itoa(int(toUserID)) + msg.Time)
	msg.ID = md5ID

	msg.IsToFromUserID = 1
	MysqlCreate <- &msg
	createUserMsgChan(msg.FromUserID)
	MessageQueue[int(msg.FromUserID)] <- &msg

	msg.IsToFromUserID = 0
	MysqlCreate <- &msg
	createUserMsgChan(msg.ToUserID)
	MessageQueue[int(msg.ToUserID)] <- &msg

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

// @Title UploadRecord
// @Description Upload media_id arr then get record url
// @Tags heal
// @Produce json
// @Param server_id body []string true "server_id"
// @Router /api/record2 [post]
// @Success 200 {object} TransformMediaIdArrToUrlResp
// @Failure 403 {object} e.ErrMsgResponse
func ConvertMediaIdArrToUrl(ctx *gin.Context) {
	var form TransformMediaIdArrToUrlReq
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(400, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	url, err := convertMediaIdArrToQiniuUrl(form.ServerID)
	if err != nil {
		ctx.JSON(500, e.ErrMsgResponse{Message: err.Error()})
		return
	}
	ctx.JSON(200, &TransformMediaIdArrToUrlResp{Url: url})
}

type TransformMediaIdArrToUrlReq struct {
	ServerID []string `json:"server_id" binding:"required"`
}

type TransformMediaIdArrToUrlResp struct {
	Url string `json:"url"`
}
