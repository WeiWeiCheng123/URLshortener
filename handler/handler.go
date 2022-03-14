package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"

	"github.com/WeiWeiCheng123/URLshortener/lib/function"
	"github.com/WeiWeiCheng123/URLshortener/model"
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

//Give a long URL, if the data format is correct, then save to DB and return a short URL.
//Otherwise, return an error and won't save to DB
func Shorten(c *gin.Context) {
	data := ShortURLForm{}
	err := c.BindJSON(&data)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	url := data.Originurl
	//Wrong URL format
	if !function.IsURL(url) {
		fmt.Println("NOT URL")
		c.String(http.StatusBadRequest, "Invalid URL")
		return
	}

	exp := data.Exp
	expTime, err := function.TimeFormater(exp)
	//Wrong Time format or time expire
	if err != nil {
		fmt.Println("ERROR TIME ", err.Error())
		c.String(http.StatusBadRequest, "Error time format or time is expired")
		return
	}

	id, err := model.Pg_Save(url, expTime)
	//Fail to save
	if err != nil {
		fmt.Println("ERROR TO SAVE ", err.Error())
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
	shortID := c.Param("shortID")
	if len(shortID) != 7 {
		fmt.Println("Length error")
		c.String(http.StatusNotFound, "This short URL is not existed or expired")
		return
	}

	if err := function.ShortID_legal(shortID); err != nil {
		fmt.Println("ShortID ilegal")
		c.String(http.StatusNotFound, "This short URL is not existed or expired")
		return
	}

	url, err := model.Redis_Load(shortID)
	if url == "NotExist" {
		fmt.Println("Not exist")
		c.String(http.StatusNotFound, "This short URL is not existed or expired")
		return
	}

	if err != nil {
		mux.RLock()
		exist, url, expireTime := model.Pg_Load(shortID)
		if !exist {
			fmt.Println("Not exist")
			model.Redis_Set_NotExist(shortID)
			mux.RUnlock()
			c.String(http.StatusNotFound, "This short URL is not existed or expired")
			return
		}

		expTime, err := function.Time_to_Taiwanzone(expireTime)
		//Wrong Time format or time expire
		if err != nil {
			fmt.Println("Expired")
			model.Redis_Set_NotExist(shortID)
			mux.RUnlock()
			c.String(http.StatusNotFound, "This short URL is not existed or expired")
			model.Pg_Del(shortID)
			return
		}

		model.Redis_Save(shortID, url, expTime)
		mux.RUnlock()
		fmt.Println("Redirect to ", url)
		c.Redirect(http.StatusFound, url)
		return
	}

	fmt.Println("Redirect to ", url)
	c.Redirect(http.StatusFound, url)
}
