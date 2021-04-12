package router

import (
	"healing2020/controller"
	"healing2020/controller/auth"
	"healing2020/controller/middleware"
	_ "healing2020/docs"
	"healing2020/pkg/tools"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"io"
	"os"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	f, _ := os.Create(tools.GetConfig("log", "location"))
	gin.DefaultWriter = io.MultiWriter(f)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	store := cookie.NewStore([]byte("healing2020"))
	r.Use(sessions.Sessions("session", store))
	if tools.IsDebug() {
		r.Use(middleware.Cors())
	}

	//开发时按群组分类，并记得按swagger格式注释
	api := r.Group("/api")
	api.Use(middleware.IdentityCheck)

	//qiniuToken
	api.GET("/qiniu/token", controller.QiniuToken)

	//个人
	api.POST("/register", controller.Register)                     //注册
	api.POST("/user/hobby", controller.NewHobby)                   //添加爱好
	api.GET("/user/hobby", controller.GetHobby)                    //获取爱好
	api.PUT("/user", controller.PutUser)                           //修改个人信息
	api.GET("/user", controller.ResponseMyPerponalPage)            //自己个人页
	api.GET("/user/others", controller.ResponseOthersPerponalPage) //他人个人页
	api.POST("/user/background", controller.ChangeBackground)      //修改用户个人背景
	api.PUT("/vod/hide_name", controller.HideName)                 //匿名

	//消息
	api.POST("/ws", controller.WsHandle)         //websocket服务
	api.POST("/broadcast", controller.Broadcast) //广播
	api.GET("/message", controller.MessagePage)  //消息首页
	api.POST("/message", controller.SendMessage) //发送消息

	//投递箱
	api.GET("/deliver/home", controller.AllDeliver)

	//歌房
	api.GET("/singsubject", controller.SingSubject)
	api.GET("/singhome", controller.SingHome)

	//抽奖
	api.GET("/lottery/allprize", controller.ALLPrize)
	api.GET("/lottery/mylottery", controller.UserLottery)
	api.GET("/lottery/money", controller.GetMoney)

	//评论
	api.GET("/getcomment", controller.GetComment)
	api.POST("/postcomment", controller.PostComment)

	//rank
	api.GET("/deliver/rank", controller.DeliverRank)
	api.GET("/songs/rank", controller.SongRank)
	api.GET("/user/rank", controller.UserRank)

	//main
	api.GET("/main/page", controller.MainMsg)

	//heal
	api.GET("/user/phone", controller.PhoneHeal)
	api.GET("/record", controller.Record)
	api.GET("/like", controller.Praise)
	api.POST("/record", controller.RecordHeal)
	api.POST("/vod", controller.VodPost)

	//swagger
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//test
	api.GET("/initest", controller.Test)

	//god view

	//login
	r.GET("/auth/fake/:id", auth.FakeLogin)

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
