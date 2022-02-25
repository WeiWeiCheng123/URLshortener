package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"

	"github.com/WeiWeiCheng123/URLshortener/model"
	"github.com/WeiWeiCheng123/URLshortener/pkg/function"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

var rdb *redis.Pool
var pdb *sql.DB
var mux sync.RWMutex

type ShortURLForm struct {
	Originurl string `json:"url"`
	Exp       string `json:"expireAt"`
}

//Connect to postgres and redis
func Init(postgres_db *sql.DB, redis_db *redis.Pool) {
	pdb = postgres_db
	rdb = redis_db
}

//Give a long URL, if the data format is correct, then save to DB and return a short URL.
//Otherwise, return an error and won't save to DB
func Shorten(c *gin.Context) {
	data := ShortURLForm{}
	err := c.BindJSON(&data)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	url := data.Originurl
	exp := data.Exp

	//Wrong URL format
	if !function.IsURL(url) {
		fmt.Println("NOT URL")
		c.String(http.StatusBadRequest, "Invalid URL")
		return
	}

	_, err = function.TimeFormater(exp)
	//Wrong Time format or time expire
	if err != nil {
		fmt.Println("ERROR TIME ", err.Error())
		c.String(http.StatusBadRequest, "Error time format or time is expired")
		return
	}

	mux.Lock()
	id, err := model.Pg_Save(pdb, url, exp)
	//Fail to save
	if err != nil {
		fmt.Println("ERROR TO SAVE ", err.Error())
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

	url, err := model.Redis_Load(rdb, shortURL)
	if url == "NotExist" {
		c.String(http.StatusNotFound, "This short URL is not existed or expired")
		return
	}

	if err != nil {
		mux.RLock()
		exist, _, url, expireTime := model.Pg_Load(pdb, shortURL)
		if !exist {
			model.Redis_Set_NotExist(rdb, shortURL)
			mux.RUnlock()
			c.String(http.StatusNotFound, "This short URL is not existed or expired")
			return
		}

		expTime, err := function.TimeFormater(expireTime)
		//Wrong Time format or time expire
		if err != nil {
			model.Redis_Set_NotExist(rdb, shortURL)
			mux.RUnlock()
			c.String(http.StatusNotFound, "This short URL is not existed or expired")
			model.Pg_Del(pdb, shortURL)
			return
		}

		model.Redis_Save(rdb, shortURL, url, expTime)
		mux.RUnlock()
		fmt.Println("Redirect to ", url)
		c.Redirect(http.StatusFound, url)
		return
	}

	fmt.Println("Redirect to ", url)
	c.Redirect(http.StatusFound, url)
}
