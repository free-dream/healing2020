package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"

	"github.com/gin-gonic/gin"
)

func GetComment(c *gin.Context) {
	strID := c.Query("id")
	Type := c.Query("Type")
	// com := Comment{
	// 	user_id : comment.user_id,
	// 	nickname : user.nickname,
	// 	Type : comment.Type,
	// 	song_id : comment.song_id,
	// 	deliver_id : comment.deliver_id,
	// 	content : comment.content,
	// }

	com, err := models.GetComment(strID, Type)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	c.JSON(200, com)
	return
}
