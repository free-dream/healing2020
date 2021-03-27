package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"

	"github.com/gin-gonic/gin"
)

type ToSaveBackground struct {
	Background int
}

//@Title ChangeBackground
//@Description 修改个人背景
//@Tags user
//@Produce json
//@Param json body ToSaveBackground true "修改后的个人背景"
//@Router /user/background [put]
//@Success 200 {object} e.ErrMsgResponse
//@Failure 403 {object} e.ErrMsgResponse
func ChangeBackground(c *gin.Context) {
	userInf := tools.GetUser()

	json := ToSaveBackground{}
	c.BindJSON(&json)

	err := models.UpdateUserOtherNow(userInf.ID, json.Background)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.ERROR_USER_SAVE_FAIL)})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: e.GetMsg(e.SUCCESS)})
}
