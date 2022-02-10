package handler

import (
	"github.com/gin-gonic/gin"
)

func Build() *gin.Engine {
	router := gin.Default()
	router.POST("/api/urls", shorten())
	router.GET("/{shortURL}", parse())

	return router
}

func shorten() gin.HandlerFunc {
	return nil
}

func parse() gin.HandlerFunc {
	return nil
}
