package auth

import (
    "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

    "healing2020/models/statements"
    "healing2020/pkg/setting"
    "healing2020/pkg/tools"
    "healing2020/pkg/e"

    "encoding/json"
    "encoding/base64"
    //"fmt"
    "time"
)

func Login(c *gin.Context) {

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
    keyname := "healing:token:"+token
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
