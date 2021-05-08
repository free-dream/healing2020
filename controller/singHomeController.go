package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SingHome(c *gin.Context) {
	subject := c.Query("subject")
	userIDInt, err := strconv.Atoi(subject)
	subjectID  := uint(userIDInt)
	sub, err := models.SingHome(subjectID)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, sub)
}

func PostSubject(c *gin.Context) {
	ID := c.Query("subject_id")
	Name := c.Query("name")
	Photo := c.Query("photo")
	Intro := c.Query("intro")
	err := models.PostSubject(ID, Name, Photo, Intro)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to add subject"})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: "发送歌房成功！"})
}
