package store

import (
	"database/sql"
	"fmt"

	"github.com/WeiWeiCheng123/URLshortener/function"
	_ "github.com/lib/pq"
)

type ShortURL struct {
	ShortID     string
	OriginalURL string
	ExpireTime  string
}

var db *sql.DB

func Connect_Pg() *sql.DB {
	db, err := sql.Open("postgres", "user=dcard_admin password=password123 dbname=dcard_db sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return db
}

func Pg_Save(db *sql.DB, url string, expireTime string) (string,error) {
	stmt, err := db.Prepare("INSERT INTO shortenerdb(shortid,originalurl,expiretime) VALUES($1,$2,$3)")
	defer db.Close()
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return "",err
	}
	shortID := function.Id()
	_, err = stmt.Exec(shortID, url, expireTime)
	stmt.Close()
	if err != nil {
		fmt.Println(err)
		return "",err
	}

	return shortID ,nil
}

func Pg_Load(db *sql.DB, shorturlID string) (bool, *ShortURL) {
	// If there is not exist, return false, otherwise return true
	defer db.Close()
	res := ShortURL{}
	stmt, err := db.Prepare("SELECT shortid, originalurl, expiretime FROM shortenerdb WHERE shortid = $1")
	if err != nil {
		return false, nil
	}

	err = stmt.QueryRow(shorturlID).Scan(&res.ShortID, &res.OriginalURL, &res.ExpireTime)
	stmt.Close()
	if err != nil {
		return false, nil
	}

	return true, &res
}

func Pg_Del(db *sql.DB, shorturlID uint64) error {
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM shortenerdb WHERE shortid = $1")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = stmt.Exec(shorturlID)
	stmt.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
