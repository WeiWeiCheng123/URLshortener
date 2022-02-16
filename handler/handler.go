package handler

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/WeiWeiCheng123/URLshortener/function"
	"github.com/WeiWeiCheng123/URLshortener/store"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

var rdb *redis.Pool
var pdb *sql.DB

//Connect to redis, postgres, and create a router
func Build() *gin.Engine {
	rdb = store.NewPool("127.0.0.1:6379")
	pdb = store.Connect_Pg()
	router := gin.Default()
	router.POST("/api/v1/urls", Shorten)
	router.GET("/:shortURL", Parse)
	router.Run(":8080")
	return router
}

//Give a long URL, if the data format is correct, then save to DB and return a short URL.
//Otherwise, return an error and won't save to DB
func Shorten(c *gin.Context) {
	data, _ := ioutil.ReadAll(c.Request.Body)
	postdata := string(data)
	post_split := strings.Split(postdata, ",")
	url := post_split[0][6:]
	exp := post_split[1][9 : len(post_split[1])-2]
	//Wrong URL format
	if !function.IsURL(url) {
		c.String(http.StatusBadRequest, "Invalid URL")
		return
	}

	_, err := function.TimeFormater(exp)
	//Wrong Time format or time expire
	if err != nil {
		c.String(http.StatusBadRequest, "Error time format or time is expired")
		return
	}
	id, err := store.Pg_Save(pdb, url, exp)
	//	_, err = store.Redis_Save(rdb, url, expTime)
	//Fail to save
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       id,
		"shortURL": "http://localhost:8080/" + id,
	})
}

//Give a short URL, if the URL exists, then redirect to the original URL.
//Otherwise, return an error (404) and won't redirect
func Parse(c *gin.Context) {
	shortURL := c.Param("shortURL")
	fmt.Println("shortid = ", shortURL)
	url, err := store.Redis_Load(rdb, shortURL)
	if err != nil {
		exist, _, url, expireTime := store.Pg_Load(pdb, shortURL)
		if !exist {
			c.String(http.StatusNotFound, "This short URL is not existed or expired")
			return
		}

		expTime, err := function.TimeFormater(expireTime)
		//Wrong Time format or time expire
		if err != nil {
			c.String(http.StatusNotFound, "This short URL is not existed or expired")
			return
		}

		store.Redis_Save(rdb, shortURL, url, expTime)
		fmt.Println("Redirect to ", url)
		c.Redirect(http.StatusFound, url)
		return
	}

	fmt.Println("Redirect to ", url)
	c.Redirect(http.StatusFound, url)
}
