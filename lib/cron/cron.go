package cron

import (
	"fmt"
	"time"

	"github.com/WeiWeiCheng123/URLshortener/model"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(database *gorm.DB) {
	db = database
}

// at a specific time use cron job to delete expired data
func Del_Expdata() {
	c := cron.New()
	//demo用，所以設置每當分鐘數為 5 的倍數時執行一次 (0, 5, 10, 15, 20...)
	c.AddFunc("*/5 * * * *",
		func() {
			fmt.Println("Cron Job start", time.Now())
			res := db.Delete(&model.Shortener{}, "expire_time < ?", time.Now())
			if res.Error != nil {
				fmt.Println(res.Error)
			} else {
				fmt.Println("Cron Job done, delete", res.RowsAffected, " data")
			}
		},
	)
	/*
			實際上使用，我的設計會是在每天的最冷門時段進行刪除
			在這邊我設置成凌晨三點
			*     *     *      *      *
			分    時    日     月     星期
			0-59  0-23  1-31  1-12  0-6 (週日~週六)

		c.AddFunc("0 3 * * *",
			func() {
				fmt.Println("Cron Job start", time.Now())
				res := db.Delete(&model.Shortener{}, "expire_time < ?", time.Now())
				if res.Error != nil {
					fmt.Println(res.Error)
				} else {
					fmt.Println("Cron Job done, delete", res.RowsAffected, " data")
				}
			},
		)
	*/
	c.Start()
}
