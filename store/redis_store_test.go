package store

import (
	"testing"
	"time"

	"github.com/WeiWeiCheng123/URLshortener/function"
)

func Test_Save_Pass(t *testing.T) {
	//First, save a data (correct form) to redis
	//It should not have error
	_, err := Redis_Save(function.Id(), "https://www.dcard.tw/f", time.Now().Add(5*time.Minute))
	if err != nil {
		t.Error("Error in Save")
	}
}

func Test_Load_Pass(t *testing.T) {
	//First, save a data (correct form) to redis and get the short URL
	//Then load the short URL, it should not have error (exist).
	res, err := Redis_Save(function.Id(), "https://www.dcard.tw/f", time.Now().Add(5*time.Minute))
	if err != nil {
		t.Error("Error in Save")
	}

	_, err = Redis_Load(res)
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
