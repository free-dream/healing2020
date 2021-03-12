package cron

import (
    "healing2020/models"
    "github.com/robfig/cron" 
)

func CronInit() *cron.Cron{
    c := cron.New()
    c.AddFunc("0 0 0 * *", func() {
        models.SendDeliverRank()
    })

    return c
}

