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
		token := session.Get("user")

		if startWith(rUrl, "/auth") || startWith(rUrl, "/wx") || startWith(rUrl, "/api/broadcast") {
			c.Next()
			return
		}
		if token == nil {
			if startWith(rUrl, "/api") {
				c.JSON(401, e.ErrMsgResponse{Message: "fail to authenticate"})
				c.Abort()
				return
			} else {
				redirect := c.Query("redirect")
				url := "https://healing2020.100steps.top/wx/jump2wechat?redirect=" + redirect
				c.Redirect(302, url)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func startWith(rUrl string, uri string) bool {
	if tools.IsDebug() {
		uri = "/test" + uri
	}
	if len(uri) > len(rUrl) {
		return false
	}
	rUrl = rUrl[0:len(uri)]
	return rUrl == uri
}
