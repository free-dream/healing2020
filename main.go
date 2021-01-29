package main

import (
    "healing2020/router"
    "healing2020/models"
    "healing2020/pkg/setting"
)

// @Title healing2020
// @Version 1.0
// @Description 2020治愈系

func main() {
    setting.MysqlConnTest()
    models.TableInit()
    
    routersInit := router.InitRouter()

    routersInit.Run(":8001")
}
