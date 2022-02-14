package store

import (
	"fmt"
	"testing"
)

func Test_Pg_Con(t *testing.T) {
	Connect_Pg()
}

func Test_Pg_save(t *testing.T) {
	db := Connect_Pg()
	err := Pg_Save(db, 123, "http", "2021")
	fmt.Println(err)
	if err != nil {
		t.Error("error")
	}
}
