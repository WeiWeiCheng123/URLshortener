package main

import (
	"fmt"

	"github.com/WeiWeiCheng123/URLshortener/handler"
	"github.com/WeiWeiCheng123/URLshortener/lib/config"
	"github.com/WeiWeiCheng123/URLshortener/lib/cron"
	"github.com/WeiWeiCheng123/URLshortener/lib/middleware"
	"github.com/WeiWeiCheng123/URLshortener/model"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

//Connect to postgres and redis
func init() {
	pg_connect := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		config.GetStr("DB_HOST"), config.GetStr("DB_PORT"), config.GetStr("DB_NAME"),
		config.GetStr("DB_USERNAME"), config.GetStr("DB_PASSWORD"), config.GetStr("DB_SSL_MODE"))
	pdb := model.Connect_Pg(pg_connect)
	rdb := model.NewPool(config.GetStr("REDIS_HOST"), config.GetInt("REDIS_POOL"), config.GetStr("REDIS_PASSWORD"))
	model.Init(pdb, rdb)
	middleware.Init(config.GetInt("IPLimitMax"), config.GetInt("IPLimitPeriod"))
}

//Create router
func engine() *gin.Engine {
	router := gin.Default()
	router.POST("/api/v1/urls", handler.Shorten)
	router.GET("/:shortURL", middleware.IPLimiter(), handler.Parse)
	return router
}

func main() {
	router := engine()
	cron.Del_Expdata()
	router.Run(":8080")
}


/*


請求 佛祖 神明 耶穌 祖先們 保佑
程式無Bug


⠂⠂⠂⠂⠂⠂⠂⠂▀████▀▄▄⠂⠂⠂⠂⠂⠂⠂⠂⠂⠂⠂⠂⠂⠂▄█
⠂⠂⠂⠂⠂⠂⠂⠂⠂⠂█▀░░░░▀▀▄▄▄▄▄⠂⠂⠂⠂▄▄▀▀█ 
⠂⠂⠂▄⠂⠂⠂⠂⠂⠂⠂█░░░░░░░░░░░▀▀▀▀▄░░▄▀ 
⠂▄▀░▀▄⠂⠂⠂⠂⠂⠂▀▄░░░░░░░░░░░░░░▀▄▀ 
▄▀░░░░█⠂⠂⠂⠂⠂⠂█▀░░░▄█▀▄░░░░░░▄█ 
▀▄░░░░░▀▄⠂⠂⠂█░░░░░▀██▀░░░░░██▄█ 
⠂⠂▀▄░░░░▄▀⠂█░░░▄██▄░░░▄░░▄░░▀▀░█ 
⠂⠂⠂█░░▄▀⠂⠂█░░░░▀██▀░░░░▀▀░▀▀░░▄▀
⠂⠂█░░░█⠂⠂█░░░░░░▄▄░░░░░░░░░░░▄▀ 
⠂█░░░█⠂⠂█▄▄░░░░░░░▀▀▄░░░░░░▄░█ 
⠂⠂▀▄░▄█▄█▀██▄░░▄▄░░░▄▀░░▄▀▀░░░█ 
⠂⠂⠂⠂▀███░░░░░░░░░▀▀▀░░░░▀▄░░░▄▀ 
⠂⠂⠂⠂⠂⠂▀▀█░░░░░░░░░▄░░░░░░▄▀█▀ 
⠂⠂⠂⠂⠂⠂⠂⠂▀█░░░░░▄▄▄▀░░▄▄▀▀░▄▀ 
⠂⠂⠂⠂⠂⠂⠂⠂⠂⠂▀▀▄▄▄▄▀⠂▀▀▀⠂▀▀▄▄▄▀


♂♀█▀▀▀▀▀▄▄♂♀♂▄▄▄▀▀▀▀▀▄▄▀▀▀▄ 
♂█░████░░░░░░░▀▄░░░▄▄░░░░██ 
♂▀▄░░░░▄▀▄█▀▀▄░▀▄░▀▄░░░░▄▀█ 
♂█░░░░░▀▄▀███▀░▄▀░░░░░███▀░█▄ 
█░░░░░░░░▀▄▄▄▄▀░░░░░░░░▀▄▄▀░█ 
█░░░░▄▀▀▄░░░░░░░░░█▀▀█░░░░░░█▀▄ 
♀█░░░█░░░░░░░░░░░░▀▄▀░░░░░░█▄▄▀
♀█░▄█▄░░░░▀▀▀░▄░░░░░░░░▀█░█ 
♀█░▀█▀░░░▄▀▀▀▀▄░▀▄█░░▀░░░█ 
♀▄▀░░░░█░██▄▄█░░█░░░▀▀██░█░█ 

*/