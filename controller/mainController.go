package controller

import (
	"github.com/gin-gonic/gin"
	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"
)

// @Title GetMainMsg
// @Description 首页数据
// @Tags main
// @Produce json
// @Router /api/main/page [get]
// @Param sort query string true "1综合排序2最新发布"
// @Param language query string false "language"
// @Param style query string false "style"
// @Success 200 {object} models.MainMsg
// @Failure 403 {object} e.ErrMsgResponse
func MainMsg(c *gin.Context) {
	sort := c.Query("sort")
	language := c.Query("language")
	style := c.Query("style")
	if !tools.Valid(sort, `^[12]$`) {
		c.JSON(403, e.ErrMsgResponse{Message: "Unexpected input"})
		return
	}
	status := typeValid(language, style)
	if status == 0 && status == 3 {
		c.JSON(403, e.ErrMsgResponse{Message: "Unexpected input"})
	}
	var key string
	if status == -1 {
		key = ""
	}
	if status == 1 {
		key = language
	}
	if status == 2 {
		key = style
	}
	data, err := models.GetMainMsg(sort, key)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "Unexpected Data"})
		return
	}
	c.JSON(200, data)
	return
}

type SongType struct {
	Language []string
	Style    []string
}

func LoadType() SongType {
	language := []string{"国语", "英语", "日语", "粤语"}
	style := []string{"ACG", "流行", "古风", "民谣", "摇滚", "抖音热搜"}
	var songType SongType
	songType.Language = language
	songType.Style = style
	return songType
}

func typeValid(language string, style string) int {
	songType := LoadType()
	if language == "" && style == "" {
		return -1
	}
	var status int = 0
	for i := 0; i < len(songType.Language); i++ {
		if language == songType.Language[i] {
			status++
			break
		}
	}
	for i := 0; i < len(songType.Style); i++ {
		if style == songType.Style[i] {
			status = status + 2
			break
		}
	}
	return status
}
