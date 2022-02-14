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
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO shortenerdb(shortid,originalurl,expiretime) VALUES($1,$2,$3);")
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

func Pg_Load(db *sql.DB, shortID uint64) (bool, error) {
	defer db.Close()

	stmt, err := db.Prepare("SELECT * FROM shortenerdb WHERE shortid = $1;")
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	res, err := stmt.Exec(shortID)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	exist, err := res.RowsAffected() // = 1 means having data
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	if exist == 1 {
		return true, nil
	}
	return false, nil
}

func Pg_Del(shortID uint64) error {
	defer db.Close()
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
