package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Build() *gin.Engine {
	router := gin.Default()
	router.POST("/api/urls", Shorten())
	router.GET("/:shortURL", Parse)
	router.Run()
	return router
}

func Shorten() gin.HandlerFunc {
	return nil
}

func Parse(c *gin.Context) {
	shortURL := c.Param("shortURL")
	fmt.Println(shortURL)
}
