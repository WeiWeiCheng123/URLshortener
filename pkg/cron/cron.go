package cron
import (
	"log"

	"github.com/robfig/cron"
)

func Dosome() {
	log.Println("Starting...")

	c := cron.New()
	c.AddFunc("5 * * * *", func() {
		log.Println("hello world")
	}) // 給物件增加定時任務
	c.Start()
}
