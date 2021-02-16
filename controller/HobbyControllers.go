package controller

import (
	"strings"

	"healing2020/models"
	"healing2020/pkg/tools"

	"github.com/gin-gonic/gin"
)

var tagNum int = 7	//标签数
//1.流行 2.古风 3.民谣 4.摇滚 5.抖音热歌 6.acg 7.其它
type Tag struct {
	Tag1 int
	Tag2 int
	Tag3 int
	Tag4 int
	Tag5 int
	Tag6 int
	Tag7 int
}

func hobbyString(json Tag) string {
	//拼接要存入数据库的字符串
	h := make([]string, tagNum)
	var n int = 0
	if json.Tag1 == 1 {
		h[n] = "1"
		n++
	}
	if json.Tag2 == 1 {
		h[n] = "2"
		n++
	}
	if json.Tag3 == 1 {
		h[n] = "3"
		n++
	}
	if json.Tag4 == 1 {
		h[n] = "4"
		n++
	}
	if json.Tag5 == 1 {
		h[n] = "5"
		n++
	}
	if json.Tag6 == 1 {
		h[n] = "6"
		n++
	}
	if json.Tag7 == 1 {
		h[n] = "7"
		n++
	}	
	h = h[:n]
	return strings.Join(h,",")
}
//@Title Hobby
//@Description 爱好选择接口
//@Tags hobby
//@Produce json
//@Router /register [post]
//@Success 200 {string} string "{"message": "xxxxx"}"
//@Failure 403 {string} string "{"error": "false"}"
func Hobby(c *gin.Context) {
	//获取json
	var json Tag
	c.BindJSON(&json)
	hobby := hobbyString(json)
	//获取redis用户信息
	userInf := tools.GetUser() 
	err := models.HobbyUpdate(hobby, userInf.ID)
	if err != nil {
		c.JSON(403, gin.H{"error": err})
	}else{
		c.JSON(200, gin.H{"message": "success"})
	}
}
