package store

import (
	"fmt"
	"testing"

	"github.com/WeiWeiCheng123/URLshortener/config"
)

func Test_Pg_save(t *testing.T) {
	pg_connect := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		config.GetStr("DB_HOST"), config.GetStr("DB_PORT"), config.GetStr("DB_NAME"),
		config.GetStr("DB_USERNAME"), config.GetStr("DB_PASSWORD"), config.GetStr("DB_SSL_MODE"))
	db := Connect_Pg(pg_connect)
	_, err := Pg_Save(db, "http", "2021")
	if err != nil {
		t.Error("error")
	}
}

func Test_Pg_load_exist(t *testing.T) {
	pg_connect := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		config.GetStr("DB_HOST"), config.GetStr("DB_PORT"), config.GetStr("DB_NAME"),
		config.GetStr("DB_USERNAME"), config.GetStr("DB_PASSWORD"), config.GetStr("DB_SSL_MODE"))
	db := Connect_Pg(pg_connect)
	s := "888"
	exist, _, u, _ := Pg_Load(db, s)
	fmt.Println(exist, u)
	if exist == true {
		t.Error("error")
	}
}

func Test_Pg_load_not_exist(t *testing.T) {
	pg_connect := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		config.GetStr("DB_HOST"), config.GetStr("DB_PORT"), config.GetStr("DB_NAME"),
		config.GetStr("DB_USERNAME"), config.GetStr("DB_PASSWORD"), config.GetStr("DB_SSL_MODE"))
	db := Connect_Pg(pg_connect)
	exist, _, _, _ := Pg_Load(db, "6")
	if exist != false {
		t.Error("error")
	}
}

func Test_Pg_del_exist(t *testing.T) {
	pg_connect := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		config.GetStr("DB_HOST"), config.GetStr("DB_PORT"), config.GetStr("DB_NAME"),
		config.GetStr("DB_USERNAME"), config.GetStr("DB_PASSWORD"), config.GetStr("DB_SSL_MODE"))
	db := Connect_Pg(pg_connect)
	err := Pg_Del(db, "999")
	if err != nil {
		t.Error("error")
	}
}

func Test_Pg_del_not_exist(t *testing.T) {
	pg_connect := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
	config.GetStr("DB_HOST"), config.GetStr("DB_PORT"), config.GetStr("DB_NAME"),
	config.GetStr("DB_USERNAME"), config.GetStr("DB_PASSWORD"), config.GetStr("DB_SSL_MODE"))
	db := Connect_Pg(pg_connect)
	err := Pg_Del(db, "123")
	if err != nil {
		t.Error("error")
	}
}
