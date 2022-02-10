package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	//"github.com/WeiWeiCheng123/URLshortener/store"
)

type PostData struct {
	Url      string `json:"url" binding:"required"`
	ExpireAt string `json:"expireAt" binding:"required"`
}

func Build() *gin.Engine {
	router := gin.Default()
	router.POST("/api/urls", Shorten)
	router.GET("/:shortURL", Parse)
	router.Run(":8080")
	return router
}

func Shorten(c *gin.Context) {
	var postdata PostData
	fmt.Println(c.BindJSON(&postdata))
	c.JSON(200, gin.H{
		"url": postdata.Url,
		"exp": postdata.ExpireAt,
	})
	return
}

func Parse(c *gin.Context) {
	shortURL := c.Param("shortURL")
	//store.Load()
	fmt.Println(shortURL)
	return
}
