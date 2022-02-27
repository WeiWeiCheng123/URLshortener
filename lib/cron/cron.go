package cron

import (
	"github.com/WeiWeiCheng123/URLshortener/model"
	"github.com/robfig/cron/v3"
)

func Del_Expdata() {
	c := cron.New()
	//Demo用，所以設置成每3分鐘進行刪除
	c.AddFunc("*/3 * * * *",
		model.Pg_Del_Exp,
	)
	/*
		實際上使用，我的設計會是在每天的最冷門時段進行刪除
		在這邊我設置成凌晨三點
		*     *     *      *      *
		分    時    日     月     星期
		0-59  0-23  1-31  1-12  0-6 (週日~週六)
		c.AddFunc("0 3 * * *",
			model.Pg_Del_Exp,
		)
	*/

	c.Start()
}
