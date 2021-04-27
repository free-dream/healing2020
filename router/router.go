package router

import (
	"encoding/gob"
	"healing2020/controller"
	"healing2020/controller/auth"
	"healing2020/controller/middleware"
	_ "healing2020/docs"
	"healing2020/pkg/tools"
	"log"

	// "time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"io"
	"os"
)

var store redis.Store

func InitRouter() *gin.Engine {
	r := gin.Default()

	f, _ := os.Create(tools.GetConfig("log", "location"))
	gin.DefaultWriter = io.MultiWriter(f)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 注册sessions组件，使用redis作为驱动，顺便注册User结构体的解码
	gob.Register(tools.RedisUser{})
	var err error
	store, err = redis.NewStore(30, "tcp", tools.GetConfig("redis", "addr"), "", []byte("100steps__"))
	if err != nil {
		log.Panicln(err.Error())
	}
	r.Use(sessions.Sessions("healing2020_session", store))

	if tools.IsDebug() {
		r.Use(middleware.Cors())
	}

	r.GET("/wx/jump2wechat", auth.JumpToWechat)
	r.GET("/wx/login", auth.DisposableLogin)
	r.POST("/wx/oauth/*redirect", auth.WechatOAuth)

	//开发时按群组分类，并记得按swagger格式注释
	api := r.Group("/api")
	api.Use(middleware.IdentityCheck())

	//qiniuToken
	api.GET("/qiniu/token", controller.QiniuToken)

	//个人
	api.POST("/register", controller.Register)                     //注册
	api.POST("/user/hobby", controller.NewHobby)                   //添加爱好
	api.GET("/user/hobby", controller.GetHobby)                    //获取爱好
	api.PUT("/user", controller.PutUser)                           //修改个人信息
	api.GET("/user", controller.ResponseMyPerponalPage)            //自己个人页
	api.GET("/user/others", controller.ResponseOthersPerponalPage) //他人个人页
	api.GET("/usermodel", controller.GetUser)                      // 获取已登录用户信息
	api.POST("/user/background", controller.ChangeBackground)      //修改用户个人背景
	api.PUT("/vod/hide_name", controller.HideName)                 //匿名

	//消息
	api.GET("/ws", controller.WsHandle)          //websocket服务
	api.POST("/broadcast", controller.Broadcast) //广播
	api.GET("/message", controller.MessagePage)  //消息首页
	api.POST("/message", controller.SendMessage) //发送消息

	//投递箱
	api.GET("/deliver/home", controller.AllDeliver)
	api.POST("/deliver/postdeliver", controller.PostDeliver)

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
