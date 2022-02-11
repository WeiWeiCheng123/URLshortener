package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Shorten_Pass(t *testing.T) {
	router := Build()
	response := "'{url:https://www.dcard.tw/f,expireAt:2022-02-20T15:04:05Z}'"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/urls", strings.NewReader(response))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "http://localhost:8080/")
}

func Test_Shorten_Fail_wrong_url(t *testing.T) {
	router := Build()
	response := "'{url:https//www.dcard.tw/f,expireAt:2022-02-20T15:04:05Z}'"  // mise ":"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/urls", strings.NewReader(response))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid URL")
}

func Test_Shorten_Fail_time_wrong_format(t *testing.T) {
	router := Build()
	response := "'{url:https://www.dcard.tw/f,expireAt:2022-02-T15:04:05Z}'"  // miss day
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/urls", strings.NewReader(response))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error time format or time is expired")
}

func Test_Shorten_Fail_time_expired(t *testing.T) {
	router := Build()
	response := "'{url:https://www.dcard.tw/f,expireAt:2022-02-10T15:04:05Z}'"  // time expired
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/urls", strings.NewReader(response))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Error time format or time is expired")
}
