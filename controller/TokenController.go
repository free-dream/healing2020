package healing2020

import (
	"fmt"
	"os"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/gin-gonic/gin"
	"healing2020/pkg/tools"
)

//@Title qiniuToken
//@Description 获取七牛的upToken
//@Tags qiniu
//@Produce json
//@Router /qiniu/token [get]
//@Success 200 {string} string "{"upToken": "xxxxx"}"
//@Failure 403 {string} string "{"err": "false"}"
func qiniuToken(c *gin.Context) {
  	//设置七牛基本信息
	accessKey := tools.GetConfig("qiniu", "accessKey")
	secretKey := tools.GetConfig("qiniu", "secretKey")
	bucket := "offcial-site"
	//获取token
	mac := qbox.NewMac(accessKey, secretKey)
	putPolicy := storage.PutPolicy{
			Scope: bucket,
	}
	upToken := putPolicy.UploadToken(mac)
	err := "false"
	//返回toekn
	if upToken != nil {
		c.JSON(200, upToken)
	}else{
		c.JSON(403, err)
	}	
}
