package controller

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"healing2020/models"
	"healing2020/models/statements"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type MyPersonalPage struct {
	UserID         uint                  `json:"user_id"`
	NickName       string                `json:"name"`
	Campus         string                `json:"school"`
	More           string                `json:"more"`
	Sex            int                   `json:"sex"`
	Setting1       int                   `json:"setting1"`
	Setting2       int                   `json:"setting2"`
	Setting3       int                   `json:"setting3"`
	Avatar         string                `json:"avatar"`
	Background     int                   `json:"background"`
	AvaBackground  []int                 `json:"avaBackground"`
	RemainHideName int                   `json:"hide_number"`
	TrueName       string                `json:"truename"`
	Phone          string                `json:"phone"`
	Vod            []models.RequestSongs `json:"requestSongs"`
	Songs          []models.Songs        `json:"Songs"`
	Praise         []models.Admire       `json:"admire"`
}

type OthersPersonalPage struct {
	UserID     uint                  `json:"user_id"`
	NickName   string                `json:"name"`
	Campus     string                `json:"school"`
	More       string                `json:"more"`
	Avatar     string                `json:"avatar"`
	Background int                   `json:"background"`
	Vod        []models.RequestSongs `json:"requestSongs"`
	Songs      []models.Songs        `json:"Songs"`
	Praise     []models.Admire       `json:"admire"`
}

func SplitAvaBackgroundtoI(avaBackground string) []int {
	avaBString := strings.Split(avaBackground, ",")
	var avaBInt []int
	for _, value := range avaBString {
		valueInt, _ := strconv.Atoi(value)
		avaBInt = append(avaBInt, valueInt)
	}
	return avaBInt
}

//综合处理各项数据获取最终返回结果
func responsePage(c *gin.Context, user statements.User, my_others string) {
	var err error
	myID := tools.GetUser(c).ID
	//初始化返回数据
	page := MyPersonalPage{
		UserID:   user.ID,
		NickName: user.NickName,
		Campus:   user.Campus,
		Sex:      user.Sex,
		More:     user.More,
		Setting1: user.Setting1,
		Setting2: user.Setting2,
		Setting3: user.Setting3,
		Avatar:   user.Avatar,
		Phone:    user.Phone,
		TrueName: user.TrueName,
	}

	// if user's setting3 == 0, back to default url
	if user.Setting3 == 0 {
		page.Avatar = tools.GetAvatarUrl(user.Sex)
	}

	//补充返回数据
	userOther, err := models.ResponseUserOther(user.ID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		c.JSON(500, e.ErrMsgResponse{Message: err.Error()})
		return
	}
	page.Background = userOther.Now
	page.AvaBackground = SplitAvaBackgroundtoI(userOther.AvaBackground)
	page.RemainHideName = userOther.RemainHideName

	page.Praise, err = models.ResponsePraise(user.ID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		c.JSON(500, e.ErrMsgResponse{Message: err.Error()})
		return
	}

	if my_others == "my" {
		page.Songs, err = models.ResponseSongs(user.ID, myID, "my")
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			c.JSON(500, e.ErrMsgResponse{Message: err.Error()})
			return
		}
		page.Vod, err = models.ResponseVod(user.ID, "my")
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			c.JSON(500, e.ErrMsgResponse{Message: err.Error()})
			return
		}
		c.JSON(200, page)
	} else if my_others == "others" {
		page.Songs, err = models.ResponseSongs(user.ID, myID, "others") //加入了匿名无数据的条件
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			c.JSON(500, err.Error())
			return
		}
		page.Vod, err = models.ResponseVod(user.ID, "others")
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			c.JSON(500, err.Error())
			return
		}
		c.JSON(200, OthersPersonalPage{
			UserID:     page.UserID,
			NickName:   page.NickName,
			Campus:     page.Campus,
			More:       page.More,
			Avatar:     page.Avatar,
			Background: page.Background,
			Vod:        page.Vod,
			Songs:      page.Songs,
			Praise:     page.Praise,
		})
	}
	return
}

//@Title ResponseMyPerponalPage
//@Description 已登录用户的个人页接口
//@Tags user
//@Produce json
//@Router /api/user [get]
//@Success 200 {object} MyPersonalPage
//@Failure 403 {object} e.ErrMsgResponse
func ResponseMyPerponalPage(c *gin.Context) {
	rUser := tools.GetUser(c)
	responsePage(c, statements.User(rUser), "my")
}

//@Title ResponseOthersPerponalPage
//@Description 其它用户的个人页接口
//@Tags user
//@Produce json
//@Router /api/user/{id} [get]
//@Success 200 {object} OthersPersonalPage
//@Failure 403 {object} e.ErrMsgResponse
func ResponseOthersPerponalPage(c *gin.Context) {
	//获取querystring并转化格式
	id := c.Query("id")
	userIDInt, err := strconv.Atoi(id)
	var userID uint = uint(userIDInt)
	if err != nil {
		fmt.Println("字符串转换成整数失败")
		c.JSON(403, e.ErrMsgResponse{Message: "获取qs参数错误"})
		return
	}
	//查询id对应用户信息
	user, err := models.ResponseUser(userID)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, e.ErrMsgResponse{Message: err.Error()})
		return
	}
	responsePage(c, user, "others")
}

type VodID struct {
	VodID uint `json:"VodID"`
}

//@Title HideName
//@Description 匿名
//@Tags user
//@Produce json
//@Router /api/vod/hide_name [put]
//@Success 200 {object} e.ErrMsgResponse
//@Failure 403 {object} e.ErrMsgResponse
func HideName(c *gin.Context) {
	jsonInf := VodID{}
	c.BindJSON(&jsonInf)

	userID := tools.GetUser(c).ID
	userOther, err := models.ResponseUserOther(userID)
	if err != nil {
		log.Println(err)
		c.JSON(500, e.ErrMsgResponse{Message: err.Error()})
		return
	}
	if userOther.RemainHideName <= 0 {
		c.JSON(403, e.ErrMsgResponse{Message: "已无剩余匿名次数！"})
		return
	}

	err = models.HideName(jsonInf.VodID, userID)
	if err != nil {
		log.Println(err)
		c.JSON(500, e.ErrMsgResponse{Message: err.Error()})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: e.GetMsg(e.SUCCESS)})
}
