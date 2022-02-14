package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Connect_Pg() *sql.DB {
	db, err := sql.Open("postgres", "user=postgres password=password dbname=shortenerDB sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return db
}

func Pg_Save(db *sql.DB, shortID uint64, url string, expireTime string) error {
	stmt, err := db.Prepare("INSERT INTO shortenerDB(shortID, originalURL, expireTime) VALUES($1,$2,$3 RETURN uid")
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = stmt.Exec(shortID, url, expireTime)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func Pg_Load() error {
	res, err := db.Query("SELECT * FROM shortenerDB where shortID=$1")
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}

func Pg_Del(shortID uint64) error {
	stmt, err := db.Prepare("DELETE FROM shortenerDB where shortID=$1")
	if err != nil {
		fmt.Println(err)
		return err
	}
	res, err := stmt.Exec(shortID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}
