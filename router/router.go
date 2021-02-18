package router

import (
    _ "healing2020/docs"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/gin-gonic/gin"
    "healing2020/controller"
)

func InitRouter() *gin.Engine {
    r := gin.New()

    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    //开发时按群组分类，并记得按swagger格式注释
    
    //qiniuToken
    r.GET("/qiniu/token", controller.QiniuToken)

    //注册
    r.POST("/register", controller.Register)

    //添加爱好
    r.POST("/user/hobby", controller.NewHobby)

    //获取爱好
    r.GET("/user/hobby", controller.GetHobby)

    //修改个人信息
    r.PUT("/user", controller.PutUser)

    //swagger
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    //test
    r.GET("/initest",controller.Test)
    return r
}

// @Title name
// @Description 
// @Tags groupName
// @Produce json
// @Router /xx/xx/xx [get/post]
// @Params xxx formData string true "xxx"
// @Success 200 {object} structName
// @Failure xxx {object} structName
