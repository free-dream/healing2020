package main

import (
    "healing2020/router"
)

// @Title healing2020
// @Version 1.0
// @Description 2020治愈系

func main() {
    routersInit := router.InitRouter()

    routersInit.Run(":8001")
}
