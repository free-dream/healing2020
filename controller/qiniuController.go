package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

	"healing2020/pkg/e"
	"healing2020/pkg/setting"
	"healing2020/pkg/tools"
)

type Token struct {
	UpToken string `json:"uptoken"`
}

const g_bucket = "healing2020"
const g_wechat_media_download_url = "https://api.weixin.qq.com/cgi-bin/media/get/jssdk?access_token=%s&media_id=%s"

var g_qiniu_upload_config storage.Config
var g_qiniu_upload_token string
var g_qiniu_accesskey string
var g_qiniu_secretkey string

func init() {
	g_qiniu_upload_config.Zone = &storage.ZoneHuanan
	g_qiniu_upload_config.UseHTTPS = false
	g_qiniu_upload_config.UseCdnDomains = false

	g_qiniu_accesskey = tools.GetConfig("qiniu", "accessKey")
	g_qiniu_secretkey = tools.GetConfig("qiniu", "secretKey")

	go updateUploadToken()
}

//@Title qiniuToken
//@Description 获取七牛的upToken
//@Tags qiniu
//@Produce json
//@Router /api/qiniu/token [get]
//@Success 200 {object} Token
//@Failure 403 {object} e.ErrMsgResponse
func QiniuToken(c *gin.Context) {
	//返回toekn
	token := getUploadToken()
	if token != "" {
		c.JSON(200, Token{UpToken: token})
	} else {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.ERROR_AUTH_TOKEN)})
	}
}

func convertMediaIdArrToQiniuUrl(media_id_arr []string) (string, error) {
	new_name := media_id_arr[0] + string(time.Now().Unix())
	if err := downloadSpeexFromWechat(media_id_arr); err != nil {
		return "", err
	}
	if err := decodeSpeexToWav(media_id_arr); err != nil {
		return "", err
	}
	if err := concatWav(media_id_arr, new_name); err != nil {
		return "", err
	}
	if err := convertWavToMp3(new_name); err != nil {
		return "", err
	}
	if _, err := uploadMp3ToQiniu(new_name); err != nil {
		return "", err
	}
	removeTmpFiles(media_id_arr)
	return fmt.Sprintf("http://cdn.healing2020.100steps.top/%s.mp3", new_name), nil
}

func getUploadToken() string {
	return g_qiniu_upload_token
}

type WechatServerErr struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func getAccessToken() string {
	redis_cli := setting.RedisClient
	return redis_cli.Get("apiv2:wechat:accesskey").Val()
}

func downloadSpeexFromWechat(media_id_arr []string) error {
	for _, media_id := range media_id_arr {
		// download byte data from wechat server
		resp, err := http.Get(fmt.Sprintf(g_wechat_media_download_url, getAccessToken(), media_id))
		if err != nil {
			// TODO: handle nerwork error
			return err
		}
		defer resp.Body.Close()

		// check if wechat server returns a error
		if resp.Header.Get("Content-Type") == "application/json; encoding=utf-8" || resp.Header.Get("Content-Type") == "text/plain" {
			var err_resp WechatServerErr
			byte_html, _ := ioutil.ReadAll(resp.Body)
			_ = json.Unmarshal(byte_html, &err_resp)
			return errors.New(err_resp.ErrMsg)
		}

		// write data to file
		f, err := os.Create(fmt.Sprintf("./media/spx/%s.spx", media_id))
		defer f.Close()
		if err != nil {
			// TODO: handle os error
			return err
		}
		io.Copy(f, resp.Body)
	}
	return nil
}

func decodeSpeexToWav(media_id_arr []string) error {
	for _, media_id := range media_id_arr {
		cmd := exec.Command("speex_decode", fmt.Sprintf("./media/spx/%s.spx", media_id), fmt.Sprintf("./media/wav/%s.wav", media_id))
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func concatWav(media_id_arr []string, new_name string) error {
	// remove old file
	cmd := exec.Command("rm", fmt.Sprintf("./media/wav/%s.wav", new_name))
	cmd.Run()

	// concat command for ffmpeg
	i := 0
	str1 := ""
	str2 := ""
	for _, media_id := range media_id_arr {
		str1 += fmt.Sprintf("-i ./media/wav/%s.wav ", media_id)
		str2 += fmt.Sprintf("[%d:0]", i)
		i++
	}
	str2 = fmt.Sprintf("-filter_complex %s'concat'=n=%d:v=0:a=1[out] -map [out] ./media/wav/%s.wav", str2, i, new_name)
	command := str1 + str2
	arglist := strings.Split(command, " ")

	// run command
	cmd = exec.Command("ffmpeg", arglist...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func convertWavToMp3(new_name string) error {
	// remove old file
	cmd := exec.Command("rm", fmt.Sprintf("./media/mp3/%s.mp3", new_name))
	cmd.Run()

	cmd = exec.Command("ffmpeg", "-i", fmt.Sprintf("./media/wav/%s.wav", new_name), fmt.Sprintf("./media/mp3/%s.mp3", new_name))
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func uploadMp3ToQiniu(new_name string) (storage.PutRet, error) {
	formUploader := storage.NewFormUploader(&g_qiniu_upload_config)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{},
	}
	if err := formUploader.PutFile(context.Background(), &ret, g_qiniu_upload_token, new_name+".mp3", fmt.Sprintf("./media/mp3/%s.mp3", new_name), &putExtra); err != nil {
		return storage.PutRet{}, err
	}
	return ret, nil
}

func removeTmpFiles(media_id_arr []string) {
	arglist := []string{}
	for _, media_id := range media_id_arr {
		arglist = append(arglist, fmt.Sprintf("./media/spx/%s.spx", media_id), fmt.Sprintf("./media/wav/%s.wav", media_id))
	}
	arglist = append(arglist, fmt.Sprintf("./media/wav/%s_concated.wav", media_id_arr[0]))
	cmd := exec.Command("rm", arglist...)
	_ = cmd.Run()
}

func updateUploadToken() {
	//获取token
	for {
		mac := qbox.NewMac(g_qiniu_accesskey, g_qiniu_secretkey)
		putPolicy := storage.PutPolicy{
			Scope: g_bucket,
		}
		g_qiniu_upload_token = putPolicy.UploadToken(mac)
		time.Sleep(3 * time.Minute)
	}
}
