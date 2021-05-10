package auth

import (
	"fmt"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"healing2020/models"
	"healing2020/models/statements"
	"healing2020/pkg/e"
	"healing2020/pkg/setting"
	"healing2020/pkg/tools"

	"encoding/base64"
	"encoding/gob"
	"encoding/json"

	//"fmt"
	"time"
)

var loginToken map[string]string

func init() {
	gob.Register(tools.RedisUser{})
	gob.Register(statements.User{})
	loginToken = make(map[string]string)
}

type LoginForm struct {
	OpenId   string `json:"openid"`
	NickName string `json:"nickname"`
	Avatar   string `json:"headimgurl"`
	Sex      int    `json:"sex"`
}

func Login(c *gin.Context) {
	// get the data from apiv3
	aes_key := tools.GetConfig("wechat", "")
	iv_key := tools.GetConfig("wechat", "")
	body := c.PostForm("body")
	data := tools.CFDDecrypter(aes_key, iv_key, body)

	if data == "" {
		c.JSON(403, e.ErrMsgResponse{Message: "unable to get the user info"})
		return
	}
	var loginForm LoginForm
	json.Unmarshal([]byte(data), loginForm)

	if loginForm.OpenId == "" {
		c.JSON(400, e.ErrMsgResponse{Message: "unable to get user's openid"})
		return
	}

	// save user info -- update or create -- set in redis
	random := tools.GetRandomString(16)
	token := base64.StdEncoding.EncodeToString(random)
	client := setting.RedisConn()
	keyname := "healing2020:token:" + token
	client.Set(keyname, data, time.Minute*30)
	models.UpdateOrCreate(loginForm.OpenId, loginForm.NickName, loginForm.Sex, loginForm.Avatar)
	// if err != nil {
	//     c.JSON(403,e.ErrMsgResponse{Message:"failed to login"})
	//     return
	// }

	// set token into session
	session := sessions.Default(c)
	session.Set("token", token)
	session.Save()

	// redirect
	redirect := c.Query("redirect")
	url := "https://healing2020.100steps.top/auth/check?redirect=" + redirect
	c.Redirect(302, url)
	return
}

type LoginStatus struct {
	Message string `json:"message"`
}

// @Title FakeLogin
// @Description 假登录接口
// @Tags login
// @Produce json
// @Router /fake [get]
// @Param id path string true "user id"
// @Param redirect query string false "redirect url"
// @Success 200 {object} LoginStatus
// @Failure 403 {object} e.ErrMsgResponse
func FakeLogin(c *gin.Context) {
	id := c.Param("id")
	redirect := c.Query("redirect")

	db := setting.MysqlConn()
	var user statements.User
	var redisUser tools.RedisUser
	result := db.Model(&statements.User{}).Where("id=?", id).First(&user)
	if result.Error != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "id maybe wrong"})
		return
	}

	tmp, _ := json.Marshal(user)
	json.Unmarshal(tmp, &redisUser)

	session := sessions.Default(c)
	session.Clear()
	session.Set("user", redisUser)
	if err := session.Save(); err != nil {
		c.JSON(500, e.ErrMsgResponse{Message: err.Error()})
		return
	}

	loginStatus := LoginStatus{
		Message: "ok",
	}

	if redirect != "" {
		c.Redirect(302, redirect)
		return
	}
	c.JSON(200, loginStatus)
}

func Jump(c *gin.Context) {
	redirect := c.Query("redirect")
	rawUrl := "https://healing2020.100steps.top/auth/login?redirect=" + redirect
	redirectUrl := base64.StdEncoding.EncodeToString([]byte(rawUrl))
	apiv3Url := "https://apiv3.100steps.top/api/bbtwoa/oauth/" + redirectUrl

	appid := tools.GetConfig("wechat", "")
	wechatUrl := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + appid + "&redirect_uri=" + apiv3Url + "&response_type=code&scope=snsapi_userinfo#wechat_redirect"

	c.Redirect(302, wechatUrl)
}

// 微信授权起点在这个接口，这里会重定向到微信服务器
func JumpToWechat(ctx *gin.Context) {
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
func WechatOAuth(ctx *gin.Context) {
	body := ctx.PostForm("body")
	user := &WechatUser{}
	json.Unmarshal([]byte(body), user)
	if user.OpenID == "" {
		ctx.JSON(403, e.ErrMsgResponse{Message: "decoding userdata failed"})
		return
	}
	loginToken[user.OpenID] = body
	ctx.String(200, fmt.Sprintf("https://healing2020.100steps.top/wx/login?token=%s&redirect=%s", user.OpenID, ctx.Param("redirect")[1:]))
}

// apiv3通过一次性登陆地址重定向到此处，完成登录流程
func DisposableLogin(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" || loginToken[token] == "" {
		ctx.JSON(401, &e.ErrMsgResponse{Message: e.GetMsg(401)})
		return
	}
	wechatUser := &WechatUser{}
	json.Unmarshal([]byte(loginToken[token]), wechatUser)
	models.UpdateOrCreate(wechatUser.OpenID, wechatUser.Nickname, wechatUser.Sex, wechatUser.HeadImgUrl)
	db := setting.MysqlConn()
	var redisUser tools.RedisUser
	var user statements.User
	result := db.Model(&statements.User{}).Where("open_id=?", wechatUser.OpenID).First(&user)
	if result.Error != nil {
		ctx.JSON(404, e.ErrMsgResponse{Message: "user not exists"})
		return
	}
	tmp, _ := json.Marshal(user)
	json.Unmarshal(tmp, &redisUser)

	// 生成session
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	option := sessions.Options{
		MaxAge: 3600,
	}
	session.Options(option)
	session.Set("user", redisUser)
	session.Save()
	redis_cli := setting.RedisClient

	// 加积分并且记录该用户本日登陆
	t := time.Now()
	t_zero := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	t_to_tomorrow := 24*60*60 - (t.Unix() - t_zero)
	logined := !redis_cli.SetNX(fmt.Sprintf("healing2020:logined_user:%d", user.ID), 0, time.Duration(t_to_tomorrow)*time.Second).Val()
	if !logined {
		models.FinishTask("1", user.ID)
	}

	redirectUrl := ctx.Query("redirect")
	if redirectUrl == "" {
		if tools.IsDebug() {
			ctx.Redirect(302, "https://healing2020.100steps.top/testfront/")
		} else {
			ctx.Redirect(302, "https://healing2020.100steps.top")
		}
	} else {
		ctx.Redirect(302, redirectUrl)
	}
}
