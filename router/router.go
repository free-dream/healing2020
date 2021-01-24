package router

import (
    "github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
    r := gin.New()

    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    //here to add groups
    
    //here to add routers

    //here to control the middleware

    return r
}
