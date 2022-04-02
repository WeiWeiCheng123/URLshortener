package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var input struct {
	URL string `json:"url"`
	Exp string `json:"expireAt"`
}

func Test_Shorten_Pass(t *testing.T) {
	// send a correct request
	// it should return 200 (redirect) and a shortURL content
	router := engine()

	TestTime := time.Now().Add(10 * time.Minute).Format("2006-01-02T15:04:05Z")
	input.URL = "https://www.dcard.tw/f"
	input.Exp = TestTime
	body, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "http://localhost:8080/")
}

func Test_Shorten_Fail_wrong_url(t *testing.T) {
	// send a wrong request (wrong url)
	// it should return 400 and Invalid url
	router := engine()

	TestTime := time.Now().Add(10 * time.Minute).Format("2006-01-02T15:04:05Z")
	input.URL = "https//www.dcard.tw/f"
	input.Exp = TestTime
	body, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid URL")
}

func Test_Shorten_Fail_time_wrong_format(t *testing.T) {
	// send a wrong request (wrong time format)
	// it should return 400 and Error time format or time is expired
	router := engine()

	TestTime := "2022-02-T15:04:05Z"
	input.URL = "https://www.dcard.tw/f"
	input.Exp = TestTime
	body, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error time format or time is expired")
}

func Test_Shorten_Fail_time_expired(t *testing.T) {
	// send a wrong request (expired time)
	// it should return 400 and Error time format or time is expired
	router := engine()

	TestTime := time.Now().Add(-10 * time.Minute).Format("2006-01-02T15:04:05Z")
	input.URL = "https://www.dcard.tw/f"
	input.Exp = TestTime
	body, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error time format or time is expired")
}
