package handler

import (
	"fmt"
	"net/http"

	"github.com/WeiWeiCheng123/URLshortener/model"
	"github.com/WeiWeiCheng123/URLshortener/pkg/function"
	"github.com/gin-gonic/gin"
)

func Shortentest(c *gin.Context) {
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

	id, err := model.Pg_Save(pdb, url, exp)
	//Fail to save
	if err != nil {
		fmt.Println("ERROR TO SAVE ", err.Error())
		mux.Unlock()
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       id,
		"shortURL": "http://localhost:8080/" + id,
	})
}
