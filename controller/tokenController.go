package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"

	"healing2020/pkg/e"
	"healing2020/pkg/tools"
)

type Token struct {
	UpToken string `json:"uptoken"`
}

//@Title qiniuToken
//@Description 获取七牛的upToken
//@Tags qiniu
//@Produce json
//@Router /api/qiniu/token [get]
//@Success 200 {object} Token
//@Failure 403 {object} e.ErrMsgResponse
func QiniuToken(c *gin.Context) {
	//设置七牛基本信息
	accessKey := tools.GetConfig("qiniu", "accessKey")
	secretKey := tools.GetConfig("qiniu", "secretKey")
	bucket := "offcial-site"
	//获取token
	mac := qbox.NewMac(accessKey, secretKey)
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	token := putPolicy.UploadToken(mac)
	//返回toekn
	if token != "" {
		c.JSON(200, Token{UpToken: token})
	} else {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.ERROR_AUTH_TOKEN)})
	}
}
