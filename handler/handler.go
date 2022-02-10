package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	//"github.com/WeiWeiCheng123/URLshortener/store"
)

func Build() *gin.Engine {
	router := gin.Default()
	router.POST("/api/urls", Shorten)
	router.GET("/:shortURL", Parse)
	router.Run(":8080")
	return router
}

func Shorten(c *gin.Context) {
	url := c.PostForm("url")
	expireAt := c.PostForm("time")
	fmt.Println("url: ", url, "\n exp:", expireAt)
}

func Parse(c *gin.Context) {
	shortURL := c.Param("shortURL")
	//store.Load()
	fmt.Println(shortURL)
}
