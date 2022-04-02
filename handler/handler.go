package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/WeiWeiCheng123/URLshortener/lib/constant"
	"github.com/WeiWeiCheng123/URLshortener/lib/function"
	"github.com/WeiWeiCheng123/URLshortener/model"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

// used for POST method
func Shorten(c *gin.Context) {
	var input struct {
		URL string `json:"url"`
		Exp string `json:"expireAt"`
	}
	data := model.Shortener{}
	err := c.BindJSON(&input)
	if err != nil {
		sendErr(c, http.StatusBadRequest, err.Error())
		return
	}

	data.OriginalUrl = input.URL
	// wrong URL format
	if !function.IsURL(data.OriginalUrl) {
		sendErr(c, http.StatusBadRequest, "invalid URL")
		return
	}

	exp := input.Exp
	data.ExpireTime, err = function.TimeFormater(exp)
	// wrong Time format or time expired
	if err != nil {
		sendErr(c, http.StatusBadRequest, "error time format or time is expired")
		return
	}
	db := c.MustGet(constant.DB).(*gorm.DB)
	shortID, id := function.Generator()
	data.ShortId = id
	res := db.Create(data)
	if res.Error != nil {
		sendErr(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Set(constant.StatusCode, http.StatusOK)
	c.Set(constant.Output, map[string]interface{}{"id": shortID, "shortURL": "http://localhost:8080/" + shortID})
	c.Set(constant.Error, nil)
}

// used for redirect original url from shortID
func Parse(c *gin.Context) {
	shortID := c.MustGet(constant.ShortID).(string)
	rdb := c.MustGet(constant.Cache).(*redis.Pool)
	connections := rdb.Get()
	defer connections.Close()

	URL, err := redis.String(connections.Do("GET", shortID))
	if URL == "NotExist" {
		sendErr(c, http.StatusNotFound, "this shortid is not existed or expired")
		return
	}
	// data is not in the redis
	if err != nil {
		data := model.Shortener{}
		id := function.Decode(shortID)
		db := c.MustGet(constant.DB).(*gorm.DB)
		res := db.Find(&data, id)
		if res.Error != nil {
			sendErr(c, http.StatusInternalServerError, err.Error())
			return
		}
		// shortID not existed
		if data.ShortId == 0 {
			_, err = connections.Do("SETEX", shortID, 150, "NotExist")
			if err != nil {
				sendErr(c, http.StatusInternalServerError, err.Error())
				return
			}

			sendErr(c, http.StatusNotFound, "this shortid is not existed or expired")
			return
		}

		data.ExpireTime, err = function.Time_to_Taiwanzone(data.ExpireTime)
		// data expired
		if err != nil {
			_, err = connections.Do("SETEX", shortID, 150, "NotExist")
			if err != nil {
				sendErr(c, http.StatusInternalServerError, err.Error())
				return
			}

			res = db.Delete(&data, data.ShortId)
			if res.Error != nil {
				sendErr(c, http.StatusInternalServerError, err.Error())
				return
			}

			sendErr(c, http.StatusNotFound, "this shortid is not existed or expired")
			return
		}

		ttl := int(data.ExpireTime.Sub(time.Now()).Seconds())
		if ttl > 900 {
			ttl = 900
			_, err = connections.Do("SETEX", shortID, ttl, data.OriginalUrl)
			if err != nil {
				sendErr(c, http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			_, err = connections.Do("SETEX", shortID, ttl, data.OriginalUrl)
			if err != nil {
				sendErr(c, http.StatusInternalServerError, err.Error())
				return
			}
		}

		c.Set(constant.StatusCode, http.StatusFound)
		c.Set(constant.Output, data.OriginalUrl)
		c.Set(constant.Error, nil)
		return
	}

	c.Set(constant.StatusCode, http.StatusFound)
	c.Set(constant.Output, URL)
	c.Set(constant.Error, nil)

}

// send error back
func sendErr(c *gin.Context, statuscode int, err string) {
	c.Set(constant.StatusCode, statuscode)
	c.Set(constant.Error, errors.New(err))
}
