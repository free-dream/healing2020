package middleware

import (
    "github.com/gin-gonic/gin"

    "healing2020/controller/auth"
)

func IdentityCheck(c *gin.Context){
    if auth.Authenticate(c) == 1 {
        c.Next()
    }else {
        c.Abort()
        return
    }
}
