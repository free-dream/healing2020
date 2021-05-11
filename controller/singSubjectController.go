package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"
	"log"

	"github.com/gin-gonic/gin"
)

type SpecialParams struct {
	Subject_id string `json:"subject_id" binding:"required"`
	Song       string `json:"song" binding:"required"`
	// User_id    string `json:"user_id" binding:"required"`
	Record     string `json:"record" binding:"required"`
}

func SingSubject(c *gin.Context) {
	sub, _ := models.SingSubject()
	c.JSON(200, sub)
}

func PostSpecial(c *gin.Context) {
	var params SpecialParams
	userInf := tools.GetUser(c)

	if err := c.ShouldBind(&params); err != nil {
		c.JSON(400, e.ErrMsgResponse{Message: "Uncomplete params"})
		return
	}
	err := models.PostSpecial(params.Subject_id, params.Song, userInf.ID, params.Record)
	if err != nil {
		log.Println(err)
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to add special"})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: "发送歌曲成功！"})
}
