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
	"github.com/gomodule/redigo/redis"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// connect to postgres and redis
func init() {
	pg_connect := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		config.GetStr("DB_HOST"), config.GetStr("DB_PORT"), config.GetStr("DB_NAME"),
		config.GetStr("DB_USERNAME"), config.GetStr("DB_PASSWORD"), config.GetStr("DB_SSL_MODE"))

	conn, err := gorm.Open(postgres.Open(pg_connect), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	db, err := conn.DB()
	if err != nil {
		log.Panic(err)
	}

	db.SetMaxIdleConns(config.GetInt("DB_MaxIdleConns"))
	db.SetMaxOpenConns(config.GetInt("DB_SetMaxOpenConns"))
	db.SetConnMaxLifetime(5 * time.Minute)

	rdb := &redis.Pool{
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

	cron.Init(conn)
	middleware.Init(conn, rdb)
	middleware.Init_ip(config.GetInt("IPLimitMax"), config.GetInt("IPLimitPeriod"))
}

// create router
func engine() *gin.Engine {
	router := gin.Default()
	router.POST("/api/v1/urls", middleware.Plain(), handler.Shorten)
	router.GET("/:shortID", middleware.IPLimiter(), middleware.TX(), handler.Parse)
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


⠂⠂⠂⠂⠂⠂⠂⠂⠂⠂▀████▀▄▄⠂⠂⠂⠂⠂⠂⠂⠂▄█
⠂⠂⠂⠂⠂⠂⠂⠂⠂⠂█▀░░░░▀▀▄▄▄▄▄⠂⠂▄▄▀▀█
⠂⠂⠂▄⠂⠂⠂⠂⠂⠂█░░░░░░░░░░░▀▀▀▀▄░░▄▀
⠂▄▀░▀▄⠂⠂⠂⠂⠂⠂▀▄░░░░░░░░░░░░░░▀▄▀
▄▀░░░░█⠂⠂⠂⠂⠂⠂█▀░░░▄█▀▄░░░░░░▄█
▀▄░░░░░▀▄⠂⠂⠂█░░░░░▀██▀░░░░░██▄█
⠂▀▄░░░░▄▀⠂⠂█░░░▄██▄░░░▄░░▄░░▀▀░█
⠂⠂⠂█░░▄▀⠂⠂█░░░░▀██▀░░░░▀▀░▀▀░░▄▀
⠂⠂█░░░█⠂⠂█░░░░░░▄▄░░░░░░░░░░░▄▀
⠂⠂█░░░█⠂⠂█▄▄░░░░░░░▀▀▄░░░░░░▄░█
⠂⠂▀▄░▄█▄█▀██▄░░▄▄░░░▄▀░░▄▀▀░░░█
⠂⠂⠂⠂▀███░░░░░░░░░▀▀▀░░░░▀▄░░░▄▀
⠂⠂⠂⠂⠂⠂▀▀█░░░░░░░░░▄░░░░░░▄▀█▀
⠂⠂⠂⠂⠂⠂⠂⠂⠂▀█░░░░░▄▄▄▀░░▄▄▀▀░▄▀
⠂⠂⠂⠂⠂⠂⠂⠂⠂▀▀▄▄▄▄▀⠂▀▀▀⠂▀▀▄▄▄▀


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
