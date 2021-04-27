package router

import (
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"healing2020/controller"
	"healing2020/controller/auth"
	"healing2020/controller/middleware"
	_ "healing2020/docs"
	"healing2020/models"
	"healing2020/models/statements"
	"healing2020/pkg/e"
	"healing2020/pkg/setting"
	"healing2020/pkg/tools"
	"log"
	"net/url"

	// "time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"io"
	"os"
)

var loginToken map[string]string
var store redis.Store

func InitRouter() *gin.Engine {
	loginToken = make(map[string]string)
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

	r.GET("/wx/jump2wechat", jumpToWechat)
	r.GET("/wx/login", disposableLogin)
	r.POST("/wx/oauth/*redirect", wechatOAuth)

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

// 微信授权起点在这个接口，这里会重定向到微信服务器
func jumpToWechat(ctx *gin.Context) {
	urlOfApiv3 := "https://apiv3.100steps.top"
	urlOfOAuth := "https://healing2020.100steps.top/wx/oauth/" + url.QueryEscape(url.QueryEscape(ctx.Query("redirect")))
	appid := "wx293bc6f4ee88d87d"
	// todo: redirect
	url2b64 := base64.StdEncoding.EncodeToString([]byte(urlOfOAuth))
	redirectUri := url.Values{}
	redirectUri.Set("redirect_uri", urlOfApiv3+"/api/bbtwoa/oauth/"+url2b64)
	finalRedirectUrl := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&%s&response_type=code&scope=snsapi_userinfo#wechat_redirect&test=false", appid, redirectUri.Encode())
	ctx.Redirect(302, finalRedirectUrl)
}

type WechatUser struct {
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	HeadImgUrl string `json:"headimgurl"`
	OpenID     string `json:"openid"`
}

// 微信服务器又会重定向到apiv3，apiv3访问该接口，该接口返回一次性登陆地址
func wechatOAuth(ctx *gin.Context) {
	body := ctx.PostForm("body")
	user := &WechatUser{}
	json.Unmarshal([]byte(body), user)
	if user.OpenID == "" {
		ctx.JSON(403, e.ErrMsgResponse{Message: "decoding userdata failed"})
		return
	}
	loginToken[user.OpenID] = body
	ctx.String(200, fmt.Sprintf("https://healing2020.100steps.top/wx/login?token=%s&redirect=%s", user.OpenID, ctx.Param("redirect")[1:]))
	ctx.String(200, fmt.Sprintf("http://test.scut18pie1.top/wx/login?token=%s&redirect=%s", user.OpenID, ctx.Param("redirect")[1:]))
}

// apiv3通过一次性登陆地址重定向到此处，完成登录流程
func disposableLogin(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" || loginToken[token] == "" {
		ctx.JSON(401, &e.ErrMsgResponse{Message: e.GetMsg(401)})
		return
	}
	wechatUser := &WechatUser{}
	json.Unmarshal([]byte(loginToken[token]), wechatUser)
	models.UpdateOrCreate(wechatUser.OpenID, wechatUser.Nickname, wechatUser.Sex, wechatUser.HeadImgUrl)
	db := setting.MysqlConn()
	defer db.Close()
	var redisUser tools.RedisUser
	var user statements.User
	result := db.Model(&statements.User{}).Where("open_id=?", wechatUser.OpenID).First(&user)
	if result.Error != nil {
		ctx.JSON(404, e.ErrMsgResponse{Message: "user not exists"})
		return
	}
	tmp, _ := json.Marshal(user)
	json.Unmarshal(tmp, &redisUser)

	session := sessions.Default(ctx)
	session.Set("user", redisUser)
	session.Save()

	redirectUrl := ctx.Query("redirect")
	if redirectUrl == "" {
		ctx.Redirect(302, "https://healing2020.100steps.top")
	} else {
		ctx.Redirect(302, redirectUrl)
	}
}
