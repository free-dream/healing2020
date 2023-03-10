package router

import (
	"encoding/gob"
	"time"

	//"healing2020/models"
	"healing2020/controller"
	"healing2020/controller/auth"
	"healing2020/controller/middleware"
	_ "healing2020/docs"
	"healing2020/pkg/e"
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
	//"fmt"
)

var store redis.Store

func InitRouter() *gin.Engine {
	var test_prefix string

	if tools.IsDebug() {
		test_prefix = "/test"
	} else {
		test_prefix = ""
	}
	r := gin.Default()

	f, _ := os.Create(tools.GetConfig("log", "location"))
	gin.DefaultWriter = io.MultiWriter(f)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Timeout(time.Minute))

	// 注册sessions组件，使用redis作为驱动
	gob.Register(tools.RedisUser{})
	var err error
	store, err = redis.NewStore(30, "tcp", tools.GetConfig("redis", "addr"), "", []byte("__100steps__100steps__100steps__"))
	if err != nil {
		log.Panicln(err.Error())
	}
	r.Use(sessions.Sessions("healing2020_session", store))

	if tools.IsDebug() {
		r.Use(middleware.Cors())
	}

	r.GET(test_prefix+"/ping", func(ctx *gin.Context) {
		ctx.JSON(200, e.ErrMsgResponse{Message: "pong"})
		return
	})

	r.GET(test_prefix+"/wx/jump2wechat", auth.JumpToWechat)
	r.GET(test_prefix+"/wx/login", auth.DisposableLogin)
	r.POST(test_prefix+"/wx/oauth/*redirect", auth.WechatOAuth)

	//开发时按群组分类，并记得按swagger格式注释
	api := r.Group(test_prefix + "/api")
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
	api.GET("/usermodel/:id", controller.GetUser)                  //获取已登录用户信息
	api.GET("/usermodel", controller.GetUser)                      //获取已登录用户信息
	api.POST("/user/background", controller.ChangeBackground)      //修改用户个人背景
	api.PUT("/vod/hide-name", controller.HideName)                 //匿名
	api.GET("/user/remain-num", controller.GetRemainNum)           //获取用户匿名次数和点歌次数
	api.POST("/user/postbox", controller.PostPostbox)              //用户邮箱

	//消息
	api.GET("/ws", controller.WsHandle)          //websocket服务
	api.POST("/broadcast", controller.Broadcast) //广播

	//投递箱
	api.GET("/deliver/home", controller.AllDeliver)
	api.GET("/deliver/single", controller.SingleDeliver)
	api.POST("/deliver/postdeliver", controller.PostDeliver)

	//歌房
	api.GET("/singsubject", controller.SingSubject)
	api.GET("/singhome", controller.SingHome)
	api.POST("/postsubject", controller.PostSubject) //发送歌房
	api.POST("/postspecial", controller.PostSpecial) //发送歌房歌曲

	//抽奖
	api.GET("/lottery/allprize", controller.ALLPrize)
	api.GET("/lottery/mylottery", controller.UserLottery)
	api.GET("/lottery/money", controller.GetMoney)
	api.GET("/lottery/usemoney", controller.LotteryDraw) //抽奖
	api.PUT("/lottery/earnmoney", controller.EarnMoney)
	api.GET("/lottery/gettask", controller.GetTask)       //获取每日任务
	api.GET("/lottery/finishtask", controller.FinishTask) //完成任务加积分

	//评论
	api.GET("/getcomment", controller.GetComment)
	api.POST("/postcomment", controller.PostComment)

	//rank
	api.GET("/deliver/rank", controller.DeliverRank)
	api.GET("/songs/rank", controller.SongRank)
	api.GET("/user/rank", controller.UserRank)
	// api.GET("/rank/update",func (c *gin.Context){
	//     err := models.SendSongRank()
	//     fmt.Println(err)
	//     err = models.SendUserRank()
	//     fmt.Println(err)
	//     err = models.SendDeliverRank()
	//     fmt.Println(err)
	// })

	//main
	api.GET("/main/page", controller.MainMsg)
	api.GET("/main/search", controller.MainSearch)

	//heal
	api.GET("/user/phone", controller.PhoneHeal)
	api.GET("/record", controller.Record)
	api.GET("/like", controller.Praise)
	api.GET("/unlike", controller.NoPraise)
	api.POST("/record", controller.RecordHeal)
	api.POST("/record2", controller.ConvertMediaIdArrToUrl)
	api.POST("/vod", controller.VodPost)

    //finish
    api.GET("/final/msg",controller.FinishData)

	//swagger
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//test
	api.GET("/initest", controller.Test)
	api.POST("/god/postprize", controller.PostPrize)

	//god view

	//login
	if tools.IsDebug() {
		r.GET("/test/auth/fake/:id", auth.FakeLogin)
	}

	return r
}
