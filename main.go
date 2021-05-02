package main

import (
	"healing2020/models"
	"healing2020/pkg/setting"
	"healing2020/pkg/tools"
	"healing2020/router"
	"log"

	//"healing2020/controller"
	"healing2020/cron"

	"github.com/fvbock/endless"
)

// @Title healing2020
// @Version 1.0
// @Description 2020治愈系

func main() {
	setting.MysqlConnTest()
	setting.RedisConnTest()
	models.TableInit()
	if tools.IsDebug() {
		//controller.LoadTestData()
		models.SendDeliverRank()
		models.SendUserRank()
		models.SendSongRank()
		models.SendMainMsg()
	}

	c := cron.CronInit()
	go c.Start()
	defer c.Stop()

	routers := router.InitRouter()
	server := endless.NewServer(":3011", routers)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err.Error())
	}
}
