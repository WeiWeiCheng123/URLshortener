package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/WeiWeiCheng123/URLshortener/lib/function"
)

func Test_Save_Pass(t *testing.T) {
	//First, save a data (correct form) to redis
	//It should not have error
	err := Redis_Save(function.Generator(), "https://www.dcard.tw/f", time.Now().Add(5*time.Minute))
	if err != nil {
		fmt.Println(err)
		t.Error("Error in Save")
	}
}

func Test_Load_Pass(t *testing.T) {
	//First, save a data (correct form) to redis and get the short URL
	//Then load the short URL, it should not have error (exist).
	ShortURL := function.Generator()
	err := Redis_Save(ShortURL, "https://www.dcard.tw/f", time.Now().Add(5*time.Minute))
	if err != nil {
		t.Error("Error in Save")
	}

	_, err = Redis_Load(ShortURL)
	if err != nil {
		t.Error("Error in Load")
	}
}

func Test_Load_not_exist(t *testing.T) {
	//Set the short URL to a non-existsent string
	//It should return Error (not exist).
	_, err := Redis_Load("WeiWei")
	if err == nil {
		t.Error("Error in Load")
	}
}
