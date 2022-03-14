package main

import (
	"fmt"
	"log"
	"time"

	"github.com/WeiWeiCheng123/URLshortener/handler"
	"github.com/WeiWeiCheng123/URLshortener/lib/config"
	"github.com/WeiWeiCheng123/URLshortener/lib/cron"
	"github.com/WeiWeiCheng123/URLshortener/lib/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/gomodule/redigo/redis"
	_ "github.com/joho/godotenv/autoload"
)

//Connect to postgres and redis
func init() {
	pg_connect := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		config.GetStr("DB_HOST"), config.GetStr("DB_PORT"), config.GetStr("DB_NAME"),
		config.GetStr("DB_USERNAME"), config.GetStr("DB_PASSWORD"), config.GetStr("DB_SSL_MODE"))

	db, err := xorm.NewEngine("postgres", pg_connect)
	if err != nil {
		log.Panic("DB connection initialization failed", err)
	}
	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	rdb := redis.Pool{
		MaxIdle:     config.GetInt("REDIS_POOL"),
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.GetStr("REDIS_HOST"))
			if err != nil {
				return nil, err
			}
			_, err = c.Do("AUTH", config.GetStr("REDIS_PASSWORD"))
			if err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}

	middleware.Init(db, rdb)
	middleware.Init_ip(config.GetInt("IPLimitMax"), config.GetInt("IPLimitPeriod"))
}

//Create router
func engine() *gin.Engine {
	router := gin.Default()
	router.POST("/api/v1/urls/1", middleware.Plain(), handler.ShortTest)
	router.GET("/1/:shortID", middleware.IPLimiter(), middleware.TX(), handler.ParseTest)
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
