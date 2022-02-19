package handler

import (
	"database/sql"
	"fmt"

	//	"io/ioutil"
	"net/http"
	//	"strings"
	"sync"

	"github.com/WeiWeiCheng123/URLshortener/function"
	"github.com/WeiWeiCheng123/URLshortener/store"
	"github.com/gin-gonic/gin"
//	"github.com/go-playground/locales/da"
	"github.com/gomodule/redigo/redis"
)

var rdb *redis.Pool
var pdb *sql.DB
var mux sync.RWMutex

type ShortURLForm struct {
	Originurl string `json:"url"`
	Exp       string `json:"expireAt"`
}

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
	data := ShortURLForm{}
	err := c.BindJSON(&data)
	if err != nil {
		c.String(http.StatusBadRequest,err.Error())
	}
	fmt.Println(data)
	url := data.Originurl
	exp := data.Exp

	/*
	data, _ := ioutil.ReadAll(c.Request.Body)
	postdata := string(data)
	post_split := strings.Split(postdata, ",")
	url := post_split[0][6:]
	exp := post_split[1][9 : len(post_split[1])-2]
	fmt.Println(url, exp)
	*/
	
	//Wrong URL format
	if !function.IsURL(url) {
		fmt.Println("NOT URL")
		c.String(http.StatusBadRequest, "Invalid URL")
		return
	}

	_, err = function.TimeFormater(exp)
	//Wrong Time format or time expire
	if err != nil {
		fmt.Println("ERROR TIME")
		c.String(http.StatusBadRequest, "Error time format or time is expired")
		return
	}
	mux.Lock()
	id, err := store.Pg_Save(pdb, url, exp)
	//	_, err = store.Redis_Save(rdb, url, expTime)
	//Fail to save
	if err != nil {
		fmt.Println("ERROR TO SAVE")
		mux.Unlock()
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	mux.Unlock()
	c.JSON(http.StatusOK, gin.H{
		"id":       id,
		"shortURL": "http://localhost:8080/" + id,
	})
}

//Give a short URL, if the URL exists, then redirect to the original URL.
//Otherwise, return an error (404) and won't redirect
func Parse(c *gin.Context) {
	shortURL := c.Param("shortURL")
	if len(shortURL) != 11 {
		c.String(http.StatusNotFound, "This short URL is not existed or expired")
		return
	}

	url, err := store.Redis_Load(rdb, shortURL)
	if err != nil {
		mux.RLock()
		exist, _, url, expireTime := store.Pg_Load(pdb, shortURL)
		if !exist {
			mux.RUnlock()
			c.String(http.StatusNotFound, "This short URL is not existed or expired")
			return
		}

		expTime, err := function.TimeFormater(expireTime)
		//Wrong Time format or time expire
		if err != nil {
			mux.RUnlock()
			c.String(http.StatusNotFound, "This short URL is not existed or expired")
			store.Pg_Del(pdb, shortURL)
			return
		}

		store.Redis_Save(rdb, shortURL, url, expTime)
		mux.RUnlock()
		fmt.Println("Redirect to ", url)
		c.Redirect(http.StatusFound, url)
		return
	}

	fmt.Println("Redirect to ", url)
	c.Redirect(http.StatusFound, url)
}
