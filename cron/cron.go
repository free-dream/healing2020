package cron

import (
	"healing2020/models"

	"github.com/robfig/cron"
)

func CronInit() *cron.Cron {
	c := cron.New()

	c.AddFunc("0 */2 * * *", func() {
		models.AutoSyncPraise()
	})
	c.AddFunc("1 */2 * * *", func() {
		models.SendDeliverRank()
	})

	c.AddFunc("2 */2 * * *", func() {
		models.SendSongRank()
	})

	c.AddFunc("3 */2 * * *", func() {
		models.SendUserRank()
	})

	c.AddFunc("0 0 0 * *", func() {
		models.UpdateRankCount()
	})

	c.AddFunc("0 */4 * * *", func() {
		models.SendMainMsg()
	})

	c.AddFunc("0 0 0 * *", func() {
		models.UpdateTask()
	})
	c.AddFunc("0 0 0 * *", func() {
		models.UpdateRemainSingDay()
	})
	return c
}
