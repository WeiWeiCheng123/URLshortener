package store

import (
	"testing"
)

func Test_Pg_save(t *testing.T) {
	db := Connect_Pg()
	_, err := Pg_Save(db, "http", "2021")
	if err != nil {
		t.Error("error")
	}
}

func Test_Pg_load_exist(t *testing.T) {
	db := Connect_Pg()
	exist, _ := Pg_Load(db, "5OrfwXgJrmO")
	if exist != true {
		t.Error("error")
	}
}

func Test_Pg_load_not_exist(t *testing.T) {
	db := Connect_Pg()
	exist, _ := Pg_Load(db, "6")
	if exist != false {
		t.Error("error")
	}
}

func Test_Pg_del_exist(t *testing.T) {
	db := Connect_Pg()
	err := Pg_Del(db, 999)
	if err != nil {
		t.Error("error")
	}
}

func Test_Pg_del_not_exist(t *testing.T) {
	db := Connect_Pg()
	err := Pg_Del(db, 123)
	if err != nil {
		t.Error("error")
	}
}
