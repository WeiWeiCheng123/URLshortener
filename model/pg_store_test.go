package model

import (
	"fmt"
	"testing"
	"time"
)

func Test_Pg_save(t *testing.T) {
	_, err := Pg_Save("http", time.Now())
	if err != nil {
		t.Error("error")
	}
}

func Test_Pg_load_exist(t *testing.T) {
	exist, _, u, _ := Pg_Load("888")
	fmt.Println(exist, u)
	if exist == true {
		t.Error("error")
	}
}

func Test_Pg_load_not_exist(t *testing.T) {
	exist, _, _, _ := Pg_Load("6")
	if exist != false {
		t.Error("error")
	}
}

func Test_Pg_del_exist(t *testing.T) {
	err := Pg_Del("999")
	if err != nil {
		t.Error("error")
	}
}

func Test_Pg_del_not_exist(t *testing.T) {
	err := Pg_Del("123")
	if err != nil {
		t.Error("error")
	}
}
