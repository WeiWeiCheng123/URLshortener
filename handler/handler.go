package handler

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/WeiWeiCheng123/URLshortener/lib/constant"
	"github.com/WeiWeiCheng123/URLshortener/lib/function"
	"github.com/WeiWeiCheng123/URLshortener/model"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/gomodule/redigo/redis"
)

var mux sync.RWMutex

func ShortTest(c *gin.Context) {
	var input struct {
		URL string `json:"url"`
		Exp string `json:"expireAt"`
	}
	err := c.BindJSON(&input)
	if err != nil {
		c.Set(constant.StatusCode, http.StatusBadRequest)
		c.Set(constant.Error, err.Error())
		return
	}

	url := input.URL
	//Wrong URL format
	if !function.IsURL(url) {
		fmt.Println("NOT URL")
		c.Set(constant.StatusCode, http.StatusBadRequest)
		c.Set(constant.Error, "Invalid URL")
		return
	}

	exp := input.Exp
	expTime, err := function.TimeFormater(exp)
	//Wrong Time format or time expire
	if err != nil {
		fmt.Println("ERROR TIME ", err.Error())
		c.Set(constant.StatusCode, http.StatusBadRequest)
		c.Set(constant.Error, "Error time format or time is expired")
		return
	}

	db := c.MustGet(constant.DB).(*xorm.Engine)
	ShortID := function.Generator()
	q := `INSERT INTO shortener(shortid, originalurl, expiretime) VALUES($1,$2,$3)`
	_, err = db.Exec(q, ShortID, url, expTime)
	if err != nil {
		fmt.Println("ERROR TO SAVE ", err.Error())
		c.Set(constant.StatusCode, http.StatusInternalServerError)
		c.Set(constant.Error, err.Error())
		return
	}

	c.Set(constant.StatusCode, http.StatusOK)
	c.Set(constant.Output, map[string]interface{}{"id": ShortID, "shortURL": "http://localhost:8080/" + ShortID})

}

func ParseTest(c *gin.Context) {
	shortID := c.MustGet(constant.ShortID).(string)
	fmt.Println("short id", shortID)
	rdb := c.MustGet(constant.Cache).(*redis.Pool)
	connections := rdb.Get()
	defer connections.Close()

	url, err := redis.String(connections.Do("GET", shortID))
	if url == "NotExist" {
		fmt.Println("Not exist")
		c.Set(constant.StatusCode, http.StatusNotFound)
		c.Set(constant.Error, "This short URL is not existed or expired")
		return
	}

	if err != nil {
		mux.RLock()
		data := model.Shortener{}
		db := c.MustGet(constant.DB).(*xorm.Engine)
		result, err := db.Where("shortid = ?", shortID).Get(&data)
		fmt.Println("res", result)
		if err != nil {
			fmt.Println("ERROR TO LOAD ", err.Error())
			c.Set(constant.StatusCode, http.StatusInternalServerError)
			c.Set(constant.Error, err.Error())
			return
		}
		if result == false {
			fmt.Println("Not exist")
			_, err = connections.Do("SETEX", shortID, 300, "NotExist")
			if err != nil {
				fmt.Println("ERROR", err.Error())
				c.Set(constant.StatusCode, http.StatusInternalServerError)
				c.Set(constant.Error, err.Error())
				return
			}

			mux.RUnlock()
			c.Set(constant.StatusCode, http.StatusNotFound)
			c.Set(constant.Error, "This short URL is not existed or expired")
			return
		}

		expTime, err := function.Time_to_Taiwanzone(data.ExpireTime)
		//Wrong Time format or time expire
		if err != nil {
			fmt.Println("Expired")
			_, err = connections.Do("SETEX", shortID, 300, "NotExist")
			if err != nil {
				fmt.Println("ERROR", err.Error())
				c.Set(constant.StatusCode, http.StatusInternalServerError)
				c.Set(constant.Error, err.Error())
				return
			}

			_, err = db.Where("shortid = ?", shortID).Delete(&data)
			if err != nil {
				fmt.Println("ERROR", err.Error())
				c.Set(constant.StatusCode, http.StatusInternalServerError)
				c.Set(constant.Error, err.Error())
				return
			}

			mux.RUnlock()
			c.Set(constant.StatusCode, http.StatusNotFound)
			c.Set(constant.Error, "This short URL is not existed or expired")
			return
		}

		_, err = connections.Do("SETEX", shortID, int(expTime.Sub(time.Now()).Seconds()), data.OriginalUrl)
		if err != nil {
			mux.RUnlock()
			fmt.Println("ERROR TO SET ", err.Error())
			c.Set(constant.StatusCode, http.StatusInternalServerError)
			c.Set(constant.Error, err.Error())
			return
		}

		c.Set(constant.StatusCode, http.StatusFound)
		c.Set(constant.Output, data.OriginalUrl)
		fmt.Println("Redirect to ", data.OriginalUrl)
	}

	c.Set(constant.StatusCode, http.StatusFound)
	c.Set(constant.Output, url)

}

//Give a short URL, if the URL exists, then redirect to the original URL.
//Otherwise, return an error (404) and won't redirect
func Parse(c *gin.Context) {
	/*
		shortID := c.Param("shortID")
		if len(shortID) != 7 {
			fmt.Println("Length error")
			c.String(http.StatusNotFound, "This short URL is not existed or expired")
			return
		}

		if err := function.ShortID_legal(shortID); err != nil {
			fmt.Println("ShortID illegal")
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
	*/
}

//Give a long URL, if the data format is correct, then save to DB and return a short URL.
//Otherwise, return an error and won't save to DB
func Shorten(c *gin.Context) {
	/*
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
	*/
}
