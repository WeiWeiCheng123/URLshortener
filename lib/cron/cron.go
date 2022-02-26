package cron
import (
	"github.com/WeiWeiCheng123/URLshortener/model"
	"github.com/robfig/cron"
)

func Del_Expdata() {
	c := cron.New()
	c.AddFunc("0 * * * *", 
		model.Pg_Del_Exp,
	)
	c.Start()
}
