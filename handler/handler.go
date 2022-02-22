package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/WeiWeiCheng123/URLshortener/function"
	"github.com/WeiWeiCheng123/URLshortener/store"
	"github.com/gin-gonic/gin"
)

var mux sync.RWMutex

type PostURLForm struct {
	Originurl string `json:"url"`
	Exp       string `json:"expireAt"`
}

//Give a long URL, if the data format is correct, then save to DB and return a short URL.
//Otherwise, return an error and won't save to DB
func Shorten(c *gin.Context) {
	fmt.Println("IP= ", c.ClientIP())
	data := PostURLForm{}
	err := c.BindJSON(&data)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	//Wrong URL format
	if !function.IsURL(data.Originurl) {
		fmt.Println("NOT URL")
		c.String(http.StatusBadRequest, "Invalid URL")
		return
	}

	_, err = function.TimeFormater(data.Exp)
	//Wrong Time format or time expire
	if err != nil {
		fmt.Println("ERROR TIME ", err.Error())
		c.String(http.StatusBadRequest, "Error time format or time is expired")
		return
	}

	mux.Lock()
	id, err := store.Pg_Save(data.Originurl, data.Exp)
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
	fmt.Println("IP= ", c.ClientIP())
	shortURL := c.Param("shortURL")
	if len(shortURL) != 11 {
		c.String(http.StatusNotFound, "This short URL is not existed or expired")
		return
	}

	url, err := store.Redis_Load(shortURL)
	if err != nil {
		mux.RLock()
		exist, data := store.Pg_Load(shortURL)
		if !exist {
			mux.RUnlock()
			c.String(http.StatusNotFound, "This short URL is not existed or expired")
			return
		}

		expTime, err := function.TimeFormater(data.ExpireTime)
		//Wrong Time format or time expire
		if err != nil {
			mux.RUnlock()
			c.String(http.StatusNotFound, "This short URL is not existed or expired")
			store.Pg_Del(shortURL)
			return
		}

		store.Redis_Save(shortURL, data.OriginalURL, expTime)
		mux.RUnlock()
		fmt.Println("Redirect to ", data.OriginalURL)
		c.Redirect(http.StatusFound, data.OriginalURL)
		return
	}

	fmt.Println("Redirect to ", url)
	c.Redirect(http.StatusFound, url)
}
