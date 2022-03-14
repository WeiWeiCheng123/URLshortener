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

func Test_Parse_Pass(t *testing.T) {
	//First send a post request and set the expire time to 10 second after
	//Then use this shortURL, it should return 302 (redirect)
	router := engine()

	TestTime := time.Now().Add(10 * time.Minute).Format("2006-01-02T15:04:05Z")
	input.URL = "https://www.dcard.tw/f"
	input.Exp = TestTime
	body, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	shortURL := w.Body.String()[7:14]
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/"+shortURL, nil)
	router.ServeHTTP(w, req)
	fmt.Println(w.Body.String())
	assert.Contains(t, w.Body.String(), "")
	assert.Equal(t, http.StatusFound, w.Code)
}

func Test_Parse_Fail_wrong_url(t *testing.T) {
	//Use an illegal shortURL id, it should return 404
	router := engine()

	w := httptest.NewRecorder()
	shortURL := "WeiWei1"
	req, _ := http.NewRequest("GET", "/"+shortURL, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "This short URL is not existed or expired")
}

func Test_Parse_Fail_url_expired(t *testing.T) {
	//First send a post request and set the expire time to 2 second after
	//Then wait for 3 second, this shortURL should expired and return 404
	router := engine()

	TestTime := time.Now().Add(2 * time.Second).Format("2006-01-02T15:04:05Z")
	input.URL = "https://www.dcard.tw/f"
	input.Exp = TestTime
	body, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	//Wait until the data expired
	time.Sleep(3 * time.Second)

	shortURL := w.Body.String()[7:14]
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/"+shortURL, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "This short URL is not existed or expired")
}
