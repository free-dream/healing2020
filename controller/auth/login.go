package auth

import (
    "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

    "healing2020/models/statements"
    "healing2020/models"
    "healing2020/pkg/setting"
    "healing2020/pkg/tools"
    "healing2020/pkg/e"

    "encoding/json"
    "encoding/base64"
    //"fmt"
    "time"
)

type LoginForm struct {
    OpenId string `json:"openid"`
    NickName string `json:"nickname"`
    Avatar string `json:"headimgurl"`
    Sex int `json:"sex"`
}

func Login(c *gin.Context) {
    // get the data from apiv3
    aes_key := tools.GetConfig("wechat","")
    iv_key := tools.GetConfig("wechat","")
    body := c.PostForm("body")
    data :=  tools.CFDDecrypter(aes_key,iv_key,body)

    if data == "" {
        c.JSON(403,e.ErrMsgResponse{Message:"unable to get the user info"})
        return
    }
    var loginForm LoginForm
    json.Unmarshal([]byte(data),loginForm)

    if loginForm.OpenId == "" {
        c.JSON(400,e.ErrMsgResponse{Message:"unable to get user's openid"})
        return 
    }

    // save user info -- update or create -- set in redis
    random := tools.GetRandomString(16)
    token := base64.StdEncoding.EncodeToString(random)
    models.UpdateOrCreate(loginForm.OpenId,loginForm.NickName,loginForm.Sex,loginForm.Avatar,token)
    // if err != nil {
    //     c.JSON(403,e.ErrMsgResponse{Message:"failed to login"})
    //     return
    // }

    // set token into session
    session := sessions.Default(c)
    session.Set("token",token)
    session.Save()

    // redirect
    redirect := c.Query("redirect")
    url := "https://healing2020.100steps.top/auth/check?redirect="+redirect
    c.Redirect(302,url)
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
// @Param id param string true "user id"
// @Param redirect query string false "redirect url"
// @Success 200 {object} LoginStatus
// @Failure 403 {object} e.ErrMsgResponse
func FakeLogin(c *gin.Context) {
    id := c.Param("id")
    redirect := c.Query("redirect")

    db := setting.MysqlConn()
    defer db.Close()
    var user statements.User
    result := db.Model(&statements.User{}).Where("id=?",id).First(&user)
    if result.Error != nil {
        c.JSON(403,e.ErrMsgResponse{Message:"id maybe wrong"})
        return
    }

    dataByte,_ := json.Marshal(user)
    data := string(dataByte)
    random := tools.GetRandomString(16)
    token := base64.StdEncoding.EncodeToString(random)

    client := setting.RedisConn()
    defer client.Close()
    keyname := "healing2020:token:"+token
    //fmt.Println(keyname)
    client.Set(keyname,data,time.Minute*30)
    //fmt.Println(result2)

    session := sessions.Default(c)
    session.Set("token",token)
    session.Save()

    loginStatus := LoginStatus{
        Message: "ok",
    }

    if redirect != "" {
        c.Redirect(302,redirect)
        return
    }
    c.JSON(200,loginStatus)
}

func Jump(c *gin.Context) {
    redirect := c.Query("redirect")
    rawUrl := "https://healing2020.100steps.top/auth/login?redirect="+redirect
    redirectUrl := base64.StdEncoding.EncodeToString([]byte(rawUrl))
    apiv3Url := "https://apiv3.100steps.top/api/bbtwoa/oauth/"+redirectUrl
    
    appid := tools.GetConfig("wechat","")
    wechatUrl := "https://open.weixin.qq.com/connect/oauth2/authorize?appid="+appid+"&redirect_uri="+apiv3Url+"&response_type=code&scope=snsapi_userinfo#wechat_redirect"

    c.Redirect(302,wechatUrl)
}
