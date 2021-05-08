package models

// import (
//     "errors"
//     "net/http"
//     "io/ioutil"
//     "fmt"
//     //"encoding/json"
//
//     "healing2020/pkg/setting"
//     "healing2020/pkg/tools"
// )
//
// func handle() {
//
// }
//
// func downloadFromWechat(mediaId string) error{
//     // 从微信服务器获取音源
//     url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential"
//     appid := tools.GetConfig("wechat","appid")
//     secret := tools.GetConfig("wechat","secret")
//     token := getAccessToken()
//     if token == "" {
//         return errors.New("fail to get access token")
//     }
//
//     url = url + "&appid=" + appid + "&secret=" + secret
//
//     cli := &http.Client{}
//     resp,err := cli.Get(url)
//     defer resp.Body.Close()
//
//     if err != nil {
//         errors.New("fail to get media")
//     }
//
//     body,_ := ioutil.ReadAll(resp.Body)
//     fmt.Println(body)
//
//     return nil
// }
//
// func getAccessToken() string{
//     client := setting.RedisConn()
//     token := client.Get("apiv3:wechat:accesskey").String()
//     return token
// }
//
// func speexToWav() {
//     // 解码把所有录音转为wav格式
//
//     // ffmepg拼接
// }
//
// func transferToMP3() {
//     // ffmepg -i xxx.wav xxx.mp3
// }
//
// func updateToQiniu() {
//     // 上传到七牛
// }
//
// func deleteTmp() {
//     // 把没必要的mp3，wav，speex等文件删除
// }
