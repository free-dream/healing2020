package controller

import (
	"fmt"
	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"
	"strconv"

	"github.com/gin-gonic/gin"
)

//@Title MessagePage
//@Description 消息首页
//@Tags message
//@Produce json
//@Router /message [get]
//@Success 200 {object} models.ToMessagePage
//@Failure 403 {object} e.ErrMsgResponse
func MessagePage(c *gin.Context) {
	user := tools.GetUser(c)
	responseMessage, err := models.ResponseMessagePage(user.ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, responseMessage)
}

//@Title CellMessage
//@Description 用户与另一用户聊天室的具体信息
//@Tags message
//@Produce json
//@Router /message/cell [GET]
//@Param id query string true "id"
//@Success 200 {object} models.ToMessageCell
//@Failure 403 {object} e.ErrMsgResponse
func CellMessage(c *gin.Context) {
	//获取querystring并转化格式
	id := c.Query("id")
	targetIDInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("字符串转换成整数失败")
		c.JSON(403, e.ErrMsgResponse{Message: "获取qs参数错误"})
		return
	}
	var targetID uint = uint(targetIDInt)

	user := tools.GetUser(c)

	responseCellMessage, err := models.SelectCellMessage(user.ID, targetID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, responseCellMessage)
}
