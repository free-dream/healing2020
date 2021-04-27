package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"healing2020/pkg/e"
	"healing2020/pkg/setting"
	"healing2020/pkg/tools"
	//"fmt"
)

func Authenticate(c *gin.Context) {
	session := sessions.Default(c)
	token := session.Get("token")
	//fmt.Println(token)
	if tools.IsZeroValue(token) {
		c.JSON(401, e.ErrMsgResponse{Message: "Fail to authenticate"})
	}

	client := setting.RedisConn()
	defer client.Close()

	_, err := client.Get("healing2020:token:" + token.(string)).Result()
	if err != nil {
		c.JSON(401, e.ErrMsgResponse{Message: "Fail to authenticate"})
	}
	// redirect := c.Query("redirect")
	// if redirect != "" {
	//     c.Redirect(302,redirect)
	// }
}
