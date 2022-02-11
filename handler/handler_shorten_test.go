package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Shorten_Pass(t *testing.T) {
	//Send a correct request
	//it should return 200 (redirect) and a shortURL content
	router := Build()
	nowTime := time.Now().Add(10 * time.Second).Format("2006-01-02T15:04:05Z")
	response := "'{url:https://www.dcard.tw/f,expireAt:" + nowTime + "}'"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/urls", strings.NewReader(response))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "http://localhost:8080/")
}

func Test_Shorten_Fail_wrong_url(t *testing.T) {
	//Send a wrong request (wrong URL)
	//it should return 400 and Invalid URL
	router := Build()
	response := "'{url:https//www.dcard.tw/f,expireAt:2022-02-20T15:04:05Z}'" // miss ":"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/urls", strings.NewReader(response))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid URL")
}

func Test_Shorten_Fail_time_wrong_format(t *testing.T) {
	//Send a wrong request (wrong time format)
	//it should return 400 and Error time format or time is expired
	router := Build()
	response := "'{url:https://www.dcard.tw/f,expireAt:2022-02-T15:04:05Z}'" // miss day
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/urls", strings.NewReader(response))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error time format or time is expired")
}

func Test_Shorten_Fail_time_expired(t *testing.T) {
	//Send a wrong request (expired time)
	//it should return 400 and Error time format or time is expired
	router := Build()
	response := "'{url:https://www.dcard.tw/f,expireAt:2022-02-10T15:04:05Z}'" // time expired
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/urls", strings.NewReader(response))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error time format or time is expired")
}
