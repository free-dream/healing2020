package main

import (
	"fmt"
	"healing2020/controller"
	"healing2020/models"
	"healing2020/pkg/setting"
	"healing2020/pkg/tools"
	"healing2020/router"
	"io/ioutil"
	"log"
	"syscall"

	//"healing2020/controller"
	"healing2020/cron"

	"github.com/fvbock/endless"
)

// @Title healing2020
// @Version 1.0
// @Description 2020治愈系

func main() {
	defer setting.DB.Close()
	defer setting.RedisClient.Close()
	models.TableInit()
	controller.MysqltoChan()
	var port string
	if tools.IsDebug() {
		//controller.LoadTestData()
		models.SendDeliverRank()
		models.SendUserRank()
		models.SendSongRank()
		models.SendMainMsg()
		port = ":3012"
	} else {
		port = ":8001"
	}

	c := cron.CronInit()
	go c.Start()
	defer c.Stop()

	// soft restart support
	routers := router.InitRouter()
	server := endless.NewServer(port, routers)
	server.BeforeBegin = func(add string) {
		pid := syscall.Getpid()
		log.Printf("Actual pid is %d", pid)
		ioutil.WriteFile("pid", []byte(fmt.Sprintf("%d", pid)), 0777)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err.Error())
	}
}
