package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func Test_Parse_Pass(t *testing.T) {
	router := Build()
	response := "'{url:https://www.dcard.tw/f,expireAt:2022-02-20T15:04:05Z}'"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/urls", strings.NewReader(response))
	router.ServeHTTP(w, req)
	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "http://localhost:8080/")
}

func Test_Parse_Fail_wrong_url(t *testing.T) {
	router := Build()
	response := "'{url:https//www.dcard.tw/f,expireAt:2022-02-20T15:04:05Z}'" // mise ":"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/urls", strings.NewReader(response))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid URL")
}
