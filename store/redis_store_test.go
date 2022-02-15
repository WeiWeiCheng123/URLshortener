package store

import (
	"testing"
	"time"

	"github.com/WeiWeiCheng123/URLshortener/function"
)

func Test_CheckId_Exist(t *testing.T) {
	//First save a data (correct form) to redis
	//Then check this data id exists or not
	//It should exist.
	r := NewPool("127.0.0.1:6379")
	id, err := Redis_Save(r, function.Id(), "https://www.dcard.tw/f", time.Now().Add(5*time.Minute))
	if err != nil {
		t.Error("Error in Save")
	}

	id_uint, _ := function.Decode(id)
	res := CheckId(r, id_uint)
	if res != true {
		t.Error("Error in Check")
	}
}

func Test_CheckId_not_Exist(t *testing.T) {
	//Set id to a non-existsent number
	//It should return false (not exist).
	r := NewPool("127.0.0.1:6379")
	res := CheckId(r, 666)

	if res != false {
		t.Error("Error in Check")
	}
}

func Test_Save_Pass(t *testing.T) {
	//First, save a data (correct form) to redis
	//It should not have error
	r := NewPool("127.0.0.1:6379")
	_, err := Redis_Save(r, function.Id(), "https://www.dcard.tw/f", time.Now().Add(5*time.Minute))
	if err != nil {
		t.Error("Error in Save")
	}
}

func Test_Load_Pass(t *testing.T) {
	//First, save a data (correct form) to redis and get the short URL
	//Then load the short URL, it should not have error (exist).
	r := NewPool("127.0.0.1:6379")
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
	r := NewPool("127.0.0.1:6379")
	_, err := Redis_Load(r, "WeiWei")
	if err == nil {
		t.Error("Error in Load")
	}
}
