package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	//"github.com/WeiWeiCheng123/URLshortener/store"
)

type PostData struct {
	OriginalURL string `json:"OriginalURL"`
	ExpireAt    string `json:"expireAt"`
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
	c.ShouldBind(&postdata)
	fmt.Println(postdata)
	c.JSON(200, gin.H{
		"url": postdata.OriginalURL,
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
