package main

import (
	"fmt"

	"github.com/WeiWeiCheng123/URLshortener/handler"
	"github.com/WeiWeiCheng123/URLshortener/lib/config"
	"github.com/WeiWeiCheng123/URLshortener/lib/store"
	"github.com/gin-gonic/gin"
	"github.com/WeiWeiCheng123/URLshortener/lib/middleware"
	_ "github.com/joho/godotenv/autoload"
)

//Connect to postgres and redis
func init() {
	pg_connect := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		config.GetStr("DB_HOST"), config.GetStr("DB_PORT"), config.GetStr("DB_NAME"),
		config.GetStr("DB_USERNAME"), config.GetStr("DB_PASSWORD"), config.GetStr("DB_SSL_MODE"))
	pdb := store.Connect_Pg(pg_connect)
	rdb := store.NewPool(config.GetStr("REDIS_HOST"), config.GetInt("REDIS_POOL"), config.GetStr("REDIS_PASSWORD"))
	fmt.Println("pdb = ", pdb)
	fmt.Println("rdb = ", rdb)
	handler.Init(pdb, rdb)
	middleware.Init(rdb)
}

//Create router
func engine() *gin.Engine {
	router := gin.Default()
	router.POST("/api/v1/urls", handler.Shorten)
	router.GET("/:shortURL", middleware.IPLimiter, handler.Parse)
	router.GET("/1/:shortURL", middleware.IPLimiter, handler.Parse1)
	return router
}

func main() {
	router := engine()
	router.Run(":8080")
}
