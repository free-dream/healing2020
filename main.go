package main

import (
    "healing2020/router"
)

func main() {
    routersInit := router.InitRouter()

    routersInit.Run(":8001")
}
