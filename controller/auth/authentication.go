package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"healing2020/pkg/e"
)

func Authenticate(c *gin.Context) {
	session := sessions.Default(c)
	data := session.Get("user")
	if data == nil {
		c.JSON(401, e.ErrMsgResponse{Message: "Fail to authenticate"})
	}

	// redirect := c.Query("redirect")
	// if redirect != "" {
	//     c.Redirect(302,redirect)
	// }
}
