package handler

import (
	"fmt"
	"net/http"

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
	fmt.Println(c.PostForm("url"))
	var Data PostData
	c.BindJSON(&Data)
	fmt.Println(&Data)
	c.JSON(http.StatusOK, gin.H{
		"url": Data.Url,
		"exp": Data.ExpireAt,
	})
	return
}

func Parse(c *gin.Context) {
	shortURL := c.Param("shortURL")
	//store.Load()
	fmt.Println(shortURL)
	return
}
