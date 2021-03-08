package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"

	"github.com/gin-gonic/gin"
)

//@Title MessagePage
//@Description 消息首页
//@Tags message
//@Produce json
//@Router /message [get]
//@Success 200 {object} e.ErrMsgResponse
//@Failure 403 {object} models.ToMessagePage
func MessagePage(c *gin.Context) {
	user := tools.GetUser()
	responseMessage, err := models.ResponseMessagePage(user.ID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, responseMessage)
}
