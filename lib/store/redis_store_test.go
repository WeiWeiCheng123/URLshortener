package store

import (
	"fmt"
	"testing"
	"time"

	"github.com/WeiWeiCheng123/URLshortener/lib/function"
)

func Test_Save_Pass(t *testing.T) {
	//First, save a data (correct form) to redis
	//It should not have error
	r := NewPool("0.0.0.0:6379", 10, "password")
	_, err := Redis_Save(r, function.Id(), "https://www.dcard.tw/f", time.Now().Add(5*time.Minute))
	if err != nil {
		fmt.Println(err)
		t.Error("Error in Save")
	}
}

func Test_Load_Pass(t *testing.T) {
	//First, save a data (correct form) to redis and get the short URL
	//Then load the short URL, it should not have error (exist).
	r := NewPool("127.0.0.1:6379", 10, "password")
	res, err := Redis_Save(r, function.Id(), "https://www.dcard.tw/f", time.Now().Add(5*time.Minute))
	if err != nil {
		t.Error("Error in Save")
	}

	_, err = Redis_Load(r, res)
	if err != nil {
		t.Error("Error in Load")
	}
}

func Test_Load_not_exist(t *testing.T) {
	//Set the short URL to a non-existsent string
	//It should return Error (not exist).
	r := NewPool("127.0.0.1:6379", 10, "password")
	_, err := Redis_Load(r, "WeiWei")
	if err == nil {
		t.Error("Error in Load")
	}
}
