package controller

import (
	"github.com/gin-gonic/gin"
	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"
)

// @Title GetUserRank
// @Description 用户积分排行榜
// @Tags rank
// @Produce json
// @Router /api/user/rank [get]
// @Success 200 {object} []models.AllRank
// @Failure 403 {object} e.ErrMsgResponse
func AllUserRank(c *gin.Context) {
	data, err := models.GetAllUserRank()
	if err != "" {
		c.JSON(403, e.ErrMsgResponse{Message: err})
		return
	}
	c.JSON(200, data)
	return
}

// @Title GetUserRank
// @Description 用户排名
// @Tags rank
// @Produce json
// @Router /api/user/rank [get]
// @Params id query string false
// @Success 200 {object} models.UserRank
// @Failure 403 {object} e.ErrMsgResponse
func UserRank(c *gin.Context) {
	if c.Query("id") == "" {
		AllUserRank(c)
		return
	}
	id := c.Query("id")
	if tools.Valid(id, `^[0-9]+$`) == false {
		c.JSON(403, e.ErrMsgResponse{Message: "error param"})
		return
	}
	data, err := models.GetUserRank(id)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "can not get rank"})
		return
	}
	c.JSON(200, data)
	return
}

// @Title GetSongRank
// @Description 每日歌曲排行榜
// @Tags rank
// @Produce json
// @Router /api/songs/rank [get]
// @Success 200 {object} []models.AllRank
// @Failure 403 {object} e.ErrMsgResponse
func SongRank(c *gin.Context) {
	data, err := models.GetSongRank(tools.GetUser(c).ID)
	if err != "" {
		c.JSON(403, e.ErrMsgResponse{Message: err})
		return
	}
	c.JSON(200, data)
	return
}

// @Title GetDeliverRank
// @Description 投递页排行榜
// @Tags rank
// @Produce json
// @Router /api/deliver/rank [get]
// @Success 200 {object} []models.AllRank
// @Failure 403 {object} e.ErrMsgResponse
func DeliverRank(c *gin.Context) {
	data, err := models.GetDeliverRank(tools.GetUser(c).ID)
	if err != "" {
		c.JSON(403, e.ErrMsgResponse{Message: err})
		return
	}
	c.JSON(200, data)
	return
}
