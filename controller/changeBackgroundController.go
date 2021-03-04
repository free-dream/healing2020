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

func ChangeBackground(c *gin.Context) {
	userInf := tools.GetUser()
	
	json := ToSaveBackground{}
	c.BindJSON(&json)

	err := models.UpdateBackgroundNow(userInf.ID, json.Background)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.ERROR_USER_SAVE_FAIL)})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: e.GetMsg(e.SUCCESS)})
}