package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"healing2020/pkg/e"
	"healing2020/pkg/tools"
)

func IdentityCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		rUrl := c.Request.URL.Path
		session := sessions.Default(c)
		token := session.Get("token")

		if startWith(rUrl, "/auth") {
			c.Next()
		}
		if tools.IsZeroValue(token) {
			if startWith(rUrl, "/api") {
				c.JSON(401, e.ErrMsgResponse{Message: "fail to authenticate"})
				c.Abort()
				return
			} else {
				redirect := c.Query("redirect")
				url := "https://healing2020.100steps.top/auth/jump?redirect=" + redirect
				c.Redirect(302, url)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func startWith(rUrl string, uri string) bool {
	if len(uri) > len(rUrl) {
		return false
	}
	rUrl = rUrl[0:len(uri)]
	return rUrl == uri
}
