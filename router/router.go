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
