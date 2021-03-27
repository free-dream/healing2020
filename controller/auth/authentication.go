package auth

import (
    "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

    "healing2020/pkg/setting"
    "healing2020/pkg/tools"
    "healing2020/pkg/e"

    "fmt"
)

func Authenticate(c *gin.Context) int{
    session := sessions.Default(c)
    token := session.Get("token")
    fmt.Println(tools.IsZeroValue(token))
    if tools.IsZeroValue(token) {
        c.JSON(401,e.ErrMsgResponse{Message:"Fail to authenticate"})
        return 0
    }

    client := setting.RedisConn()
    defer client.Close()
    _,err := client.Get("apiv3:token:"+token.(string)).Result()
    if err != nil {
        c.JSON(401,e.ErrMsgResponse{Message:"Fail to authenticate"})
        return 0
    }
    return 1
}
