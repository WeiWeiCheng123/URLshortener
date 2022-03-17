package handler

import (
	"net/http"
	"time"

	"github.com/WeiWeiCheng123/URLshortener/lib/constant"
	"github.com/WeiWeiCheng123/URLshortener/lib/function"
	"github.com/WeiWeiCheng123/URLshortener/model"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/gomodule/redigo/redis"
)

// used for POST method
func Shorten(c *gin.Context) {
	var input struct {
		URL string `json:"url"`
		Exp string `json:"expireAt"`
	}

	if err := c.BindJSON(&input); err != nil {
		sendErr(c, http.StatusBadRequest, err.Error())
		return
	}

	URL := input.URL
	// wrong URL format
	if !function.IsURL(URL) {
		sendErr(c, http.StatusBadRequest, "invalid URL")
		return
	}

	exp := input.Exp
	expTime, err := function.TimeFormater(exp)
	// wrong Time format or time expired
	if err != nil {
		sendErr(c, http.StatusBadRequest, "error time format or time is expired")
		return
	}

	db := c.MustGet(constant.DB).(*xorm.Engine)
	shortID, id := function.Generator()
	q := `INSERT INTO shortener(short_id, original_url, expire_time) VALUES($1,$2,$3)`
	_, err = db.Exec(q, id, URL, expTime)
	if err != nil {
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
		db := c.MustGet(constant.DB).(*xorm.Engine)
		result, err := db.Where("short_id = ?", id).Get(&data)
		if err != nil {
			sendErr(c, http.StatusInternalServerError, err.Error())
			return
		}
		// shortID not existed
		if result == false {
			_, err = connections.Do("SETEX", shortID, 300, "NotExist")
			if err != nil {
				sendErr(c, http.StatusInternalServerError, err.Error())
				return
			}

			sendErr(c, http.StatusNotFound, "this shortid is not existed or expired")
			return
		}

		expTime, err := function.Time_to_Taiwanzone(data.ExpireTime)
		// data expired
		if err != nil {
			_, err = connections.Do("SETEX", shortID, 300, "NotExist")
			if err != nil {
				sendErr(c, http.StatusInternalServerError, err.Error())
				return
			}

			_, err = db.Where("short_id = ?", id).Delete(&data)
			if err != nil {
				sendErr(c, http.StatusInternalServerError, err.Error())
				return
			}

			sendErr(c, http.StatusNotFound, "this shortid is not existed or expired")
			return
		}

		_, err = connections.Do("SETEX", shortID, int(expTime.Sub(time.Now()).Seconds()), data.OriginalUrl)
		if err != nil {
			sendErr(c, http.StatusInternalServerError, err.Error())
			return
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

func sendErr(c *gin.Context, statuscode int, err string) {
	c.Set(constant.StatusCode, http.StatusInternalServerError)
	c.Set(constant.Error, err)
}
