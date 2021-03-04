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
    api := r.Group("/api")
    
    //qiniuToken
    api.GET("/qiniu/token", controller.QiniuToken)

    //注册
    api.POST("/register", controller.Register)

    //添加爱好
    api.POST("/user/hobby", controller.NewHobby)
    //获取爱好
    api.GET("/user/hobby", controller.GetHobby)

    //修改个人信息
    r.PUT("/user", controller.PutUser)
    //个人页
    api.GET("/user", controller.ResponseMyPerponalPage)

    //修改用户个人背景
    api.POST("/user/background", controller.ChangeBackground)

    //rank
    api.GET("/deliver/rank",controller.DeliverRank)
    api.GET("/songs/rank",controller.SongRank)
    api.GET("/user/rank",controller.UserRank)

    //main
    api.GET("/main/page",controller.MainMsg)

    //heal
    api.GET("/user/phone",controller.PhoneHeal)
    api.GET("/record",controller.Record)
    api.GET("/like",controller.Praise)
    api.POST("/record",controller.RecordHeal)
    api.POST("/vod",controller.VodPost)

    //swagger
    api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    //test
    api.GET("/initest",controller.Test)

    //god view

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
