package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type PostForm struct {
	Originurl string `json:"url"`
	Exp       string `json:"expireAt"`
}


func Test_Parse_Pass(t *testing.T) {
	router := engine()
	router.Run()
	TestTime := time.Now().Add(10 * time.Minute).Format("2006-01-02T15:04:05Z")
	post_data := PostForm{}
	post_data.Originurl = "https://www.dcard.tw/f"
	post_data.Exp = TestTime
	body, _ := json.Marshal(post_data)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	shortURL := w.Body.String()[7:18]
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/"+shortURL, nil)
	router.ServeHTTP(w, req)
	fmt.Println(w.Body.String())
	assert.Contains(t, w.Body.String(), "")
	assert.Equal(t, http.StatusFound, w.Code)
}
