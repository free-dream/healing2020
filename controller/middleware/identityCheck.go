package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"

    "healing2020/pkg/tools"
    "healing2020/pkg/e"
)

func IdentityCheck() gin.HandlerFunc{
    return func(c *gin.Context) {
        rUrl := c.Request.URL.Path
        session := sessions.Default(c)
        token := session.Get("token")

        if startWith(rUrl,"/auth") {
            c.Next()
        }
        if tools.IsZeroValue(token) {
            if startWith(rUrl,"/api") {
                c.JSON(401,e.ErrMsgResponse{Message:"fail to authenticate"})
                c.Abort()
                return
            }else {
                redirect := c.Query("redirect")
                url := "https://healing2020.100steps.top/auth/jump"+redirect
                c.Redirect(302,url)
                c.Abort()
                return
            }
        }
        c.Next()
    }
}

func startWith(rUrl string,uri string) bool{
    rUrl = rUrl[0:len(uri)-1]
    return rUrl == uri
}
