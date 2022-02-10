package handler

import (
	"github.com/gin-gonic/gin"
)

func Build() *gin.Engine {
	router := gin.Default()
	router.POST("/api/urls", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "short",
		})
	})
	router.GET("/{shortURL}", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "origin",
		})
	})
	router.Run()
	return router
}


func Shorten() gin.HandlerFunc {
	return nil
}

func Parse() gin.HandlerFunc {
	return nil
}
