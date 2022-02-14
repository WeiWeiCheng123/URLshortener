package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Connect_Pg() *sql.DB {
	db, err := sql.Open("postgres", "user=dcard_admin password=admin_password dbname=dcard_db sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return db
}

func Pg_Save(db *sql.DB, shortID uint64, url string, expireTime string) error {
	fmt.Println(db)
	stmt, err := db.Prepare("insert into test(shortenerdb,originalurl,expiretime) values($1,$2,$3);")
	fmt.Println(stmt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	res, err := stmt.Exec(shortID, url, expireTime)
	fmt.Println("res = ", res.)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func Pg_Load() error {
	res, err := db.Query("SELECT * FROM shortenerdb where shortid=$1")
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}

func Pg_Del(shortID uint64) error {
	stmt, err := db.Prepare("DELETE FROM shortenerdb where shortid=$1")
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
