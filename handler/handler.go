package handler

import (
	"errors"
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

func Shorten(c *gin.Context) {
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
		c.Set(constant.Error, errors.New("invalid URL"))
		return
	}

	exp := input.Exp
	expTime, err := function.TimeFormater(exp)
	//Wrong Time format or time expire
	if err != nil {
		fmt.Println("ERROR TIME ", err.Error())
		c.Set(constant.StatusCode, http.StatusBadRequest)
		c.Set(constant.Error, errors.New("error time format or time is expired"))
		return
	}

	db := c.MustGet(constant.DB).(*xorm.Engine)
	ShortID := function.Generator()
	q := `INSERT INTO shortener(short_id, original_url, expire_time) VALUES($1,$2,$3)`
	_, err = db.Exec(q, ShortID, url, expTime)
	if err != nil {
		fmt.Println("ERROR TO SAVE ", err.Error())
		c.Set(constant.StatusCode, http.StatusInternalServerError)
		c.Set(constant.Error, err.Error())
		return
	}

	c.Set(constant.StatusCode, http.StatusOK)
	c.Set(constant.Output, map[string]interface{}{"id": ShortID, "shortURL": "http://localhost:8080/" + ShortID})
	c.Set(constant.Error, nil)
}

func Parse(c *gin.Context) {
	shortID := c.MustGet(constant.ShortID).(string)
	rdb := c.MustGet(constant.Cache).(*redis.Pool)
	connections := rdb.Get()
	defer connections.Close()

	url, err := redis.String(connections.Do("GET", shortID))
	if url == "NotExist" {
		fmt.Println("Not exist")
		c.Set(constant.StatusCode, http.StatusNotFound)
		c.Set(constant.Error, errors.New("this shortid is not existed or expired"))
		return
	}

	if err != nil {
		mux.RLock()
		data := model.Shortener{}
		db := c.MustGet(constant.DB).(*xorm.Engine)
		result, err := db.Where("short_id = ?", shortID).Get(&data)
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
			c.Set(constant.Error, errors.New("this shortid is not existed or expired"))
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

			_, err = db.Where("short_id = ?", shortID).Delete(&data)
			if err != nil {
				fmt.Println("ERROR", err.Error())
				c.Set(constant.StatusCode, http.StatusInternalServerError)
				c.Set(constant.Error, err.Error())
				return
			}

			mux.RUnlock()
			c.Set(constant.StatusCode, http.StatusNotFound)
			c.Set(constant.Error, errors.New("this shortid is not existed or expired"))
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
		c.Set(constant.Error, nil)
		fmt.Println("Redirect to ", data.OriginalUrl)
		return
	}

	c.Set(constant.StatusCode, http.StatusFound)
	c.Set(constant.Output, url)
	c.Set(constant.Error, nil)

}
