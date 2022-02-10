package handler

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/WeiWeiCheng123/URLshortener/store"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client


func Build() *gin.Engine {
	rdb = store.NewClient()
	router := gin.Default()
	router.POST("/api/urls", Shorten)
	router.GET("/:shortURL", Parse)
	router.Run(":8080")
	return router
}

func Shorten(c *gin.Context) {
	data, _ := ioutil.ReadAll(c.Request.Body)
	postdata := string(data)
	post_split := strings.Split(postdata, ",")
	url := post_split[0][6:]
	exp := post_split[1][9:]
	store.Save(rdb, url, exp)
	c.JSON(200, gin.H{
		"id": "",
		"shortURL" : "",
	})
	
	return
}

func Parse(c *gin.Context) {
	shortURL := c.Param("shortURL")
	//store.Load()
	fmt.Println(shortURL)
	return
}
