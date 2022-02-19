package main

import (
	"fmt"

	"github.com/WeiWeiCheng123/URLshortener/config"
	"github.com/WeiWeiCheng123/URLshortener/handler"
	"github.com/WeiWeiCheng123/URLshortener/store"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

//Connect to postgres and redis 
func init() {
	pg_connect := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		config.GetStr("DB_HOST"), config.GetStr("DB_PORT"), config.GetStr("DB_NAME"), 
		config.GetStr("DB_USERNAME"), config.GetStr("DB_PASSWORD"), config.GetStr("DB_SSL_MODE"))
	pdb := store.Connect_Pg(pg_connect)
	rdb := store.NewPool(config.GetStr("REDIS_HOST"), config.GetInt("REDIS_POOL"), config.GetStr("REDIS_PASSWORD"))

	handler.Init(pdb, rdb)
}

//Create router
func engine() *gin.Engine {
	router := gin.Default()
	router.POST("/api/v1/urls", handler.Shorten)
	router.GET("/:shortURL", handler.Parse)
	return router
}

func main() {
	router := engine()
	router.Run(":8080")
}
