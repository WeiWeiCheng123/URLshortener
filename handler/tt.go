package handler

import (
	"fmt"
	"net/http"

	"github.com/WeiWeiCheng123/URLshortener/model"
	"github.com/WeiWeiCheng123/URLshortener/pkg/function"
	"github.com/gin-gonic/gin"
)

func Parse1(c *gin.Context) {
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
		fmt.Println("hello")
		exist, _, url, expireTime := model.Pg_Load(pdb, shortURL)
		if !exist {
			model.Redis_Set_NotExist(rdb, shortURL)
			c.String(http.StatusNotFound, "This short URL is not existed or expired")
			return
		}

		expTime, err := function.TimeFormater(expireTime)
		//Wrong Time format or time expire
		if err != nil {
			model.Redis_Set_NotExist(rdb, shortURL)
			c.String(http.StatusNotFound, "This short URL is not existed or expired")
			model.Pg_Del(pdb, shortURL)
			return
		}

		model.Redis_Save(rdb, shortURL, url, expTime)
		fmt.Println("Redirect to ", url)
		c.Redirect(http.StatusFound, url)
		return
	}

	fmt.Println("Redirect to ", url)
	c.Redirect(http.StatusFound, url)
}
