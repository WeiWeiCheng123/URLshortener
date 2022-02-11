package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/WeiWeiCheng123/URLshortener/function"
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
	exp := post_split[1][9 : len(post_split[1])-2]

	if !function.IsUrl(url) {
		c.String(400, "Invalid URL")
	}

	id, err := store.Save(rdb, url, exp)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, gin.H{
		"id":       id,
		"shortURL": "http://localhost:8080/" + id,
	})
}

func Parse(c *gin.Context) {
	shortURL := c.Param("shortURL")
	url, err := store.Load(rdb, shortURL)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	fmt.Println("Redirect to ", url)
	c.Redirect(http.StatusFound, url)
}
