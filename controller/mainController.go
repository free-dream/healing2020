package controller

import (
    "net/url"
	//"fmt"

	"github.com/gin-gonic/gin"
	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"
)

// @Title Search
// @Description 首页搜索
// @Tags main
// @Produce json
// @Router /api/main/search [get]
// @Param search query string true "search form"
// @Success 200 {object} models.SearchResp
// @Failure 403 {object} e.ErrMsgResponse
func MainSearch(c *gin.Context) {
    searchRaw := c.Query("search")
    search,_ := url.QueryUnescape(searchRaw) 
    if !tools.Valid(search,"^[\u0000-\uffff]{1,16}$") {
        c.JSON(400,e.ErrMsgResponse{Message:"unexpected params"})
        return
    }
    result := models.GetSearchResult(search)
    if result.Err != "" {
        c.JSON(500,e.ErrMsgResponse{Message:"internal error"})
        return
    }
    c.JSON(200,result)
    return
}

// @Title GetMainMsg
// @Description 首页数据
// @Tags main
// @Produce json
// @Router /api/main/page [get]
// @Param sort query string true "1综合排序2最新发布"
// @Param page query string true "页数"
// @Param language query string false "language"
// @Param style query string false "style"
// @Success 200 {object} models.MainMsg
// @Failure 403 {object} e.ErrMsgResponse
func MainMsg(c *gin.Context) {
	sort := c.Query("sort")
	language := c.Query("language")
	style := c.Query("style")
    page := c.Query("page")
	if !tools.Valid(sort, `^[01]$`) || !tools.Valid(page, `^[1-9]+[0-9]*$`) {
		c.JSON(403, e.ErrMsgResponse{Message: "Unexpected input"})
		return
	}
	status := typeValid(language, style)
	if status == 0 || status == 3 {
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
    tags := ""
    if style == "推荐" {
        tags = tools.GetUser(c).Hobby
    }
	data, err := models.GetMainMsg(page,sort, key, tags)
	if err != nil {
		//c.JSON(403, e.ErrMsgResponse{Message: "Unexpected Data"})
		c.JSON(403, e.ErrMsgResponse{Message: err.Error()})
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
	style := []string{"推荐","ACG", "流行", "古风", "民谣", "摇滚", "抖音热歌", "其他"}
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
