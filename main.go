package main

import (
    "healing2020/router"
    "healing2020/models"
    "healing2020/pkg/setting"
    "healing2020/pkg/tools"
    "healing2020/controller"
    "healing2020/cron"
)

// @Title healing2020
// @Version 1.0
// @Description 2020治愈系

func main() {
    setting.MysqlConnTest()
    setting.RedisConnTest()
    models.TableInit()
    if tools.IsDebug() {
        controller.LoadTestData()
        models.SendDeliverRank()
    }

    c := cron.CronInit()
    go c.Start()
    defer c.Stop()
    
    routersInit := router.InitRouter()

    routersInit.Run(":8001")
}
