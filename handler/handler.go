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

func Build() *gin.Engine {
	rdb = store.NewPool("127.0.0.1:6379")
	pdb = store.Connect_Pg()
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

func Parse(c *gin.Context) {
	shortURL := c.Param("shortURL")
	fmt.Println("shortid = ",shortURL)
	url, err := store.Redis_Load(rdb, shortURL)
	if err != nil {
		exist, res := store.Pg_Load(pdb, shortURL)
		fmt.Println(res, exist)
		if exist == false {
			c.String(http.StatusNotFound, "1This short URL is not existed or expired")
			return
		}
		expTime, err := function.TimeFormater(res.ExpireTime)
		//Wrong Time format or time expire
		if err != nil {
			c.String(http.StatusNotFound, "2This short URL is not existed or expired")
			return
		}
		store.Redis_Save(rdb, res.OriginalURL, expTime)
	}

	fmt.Println("Redirect to ", url)
	c.Redirect(http.StatusFound, url)
}

func tt() {
	e, res := store.Pg_Load(pdb, "bd7ujchwLuU")
	fmt.Println(e, res)
}