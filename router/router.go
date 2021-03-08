package router

import (
	"healing2020/controller"
	_ "healing2020/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

<<<<<<< HEAD
	//开发时按群组分类，并记得按swagger格式注释

	//ws
	r.GET("/ws", controller.WsHandle)
	r.POST("/broadcast", controller.Broadcast)

	//消息首页
	r.GET("/message", controller.MessagePage)

	//qiniuToken
	r.GET("/qiniu/token", controller.QiniuToken)

	//爱好
	//添加爱好
	r.POST("/user/hobby", controller.NewHobby)
	//获取爱好
	r.GET("/user/hobby", controller.GetHobby)

	//个人
	//注册
	r.POST("/register", controller.Register)
	//修改个人信息
	r.PUT("/user", controller.PutUser)
	//个人页
	r.GET("/user", controller.ResponseMyPerponalPage)
	//修改用户个人背景
	r.POST("/user/background", controller.ChangeBackground)

	//rank
	r.GET("/deliver/rank", controller.DeliverRank)
	r.GET("/songs/rank", controller.SongRank)
	r.GET("/user/rank", controller.UserRank)

	//main
	r.GET("/main/page", controller.MainMsg)

	//swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//test
	r.GET("/initest", controller.Test)

	//god view
	r.POST("/new/deliver", controller.NewDeliverRank)
=======
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
>>>>>>> fae13868e69b4bdbe56ab926c618f348017ba4e6

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
