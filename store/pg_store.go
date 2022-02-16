package store

import (
	"database/sql"

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
		panic(err)
	}

	return db
}

//Give original URL and expire time, save to Postgres.
func Pg_Save(db *sql.DB, url string, expireTime string) (string, error) {
	stmt, err := db.Prepare("INSERT INTO shortenerdb(shortid,originalurl,expiretime) VALUES($1,$2,$3)")
	if err != nil {
		return "", err
	}

	shortID := function.Id()
	_, err = stmt.Exec(shortID, url, expireTime)
	defer stmt.Close()
	if err != nil {
		return "", err
	}
	return shortID, nil
}

// If there is not exist, return false, otherwise return true
func Pg_Load(db *sql.DB, shorturlID string) (bool, string, string, string) {
	stmt, err := db.Prepare("SELECT shortid, originalurl, expiretime FROM shortenerdb WHERE shortid = $1")
	if err != nil {
		return false, "", "", ""
	}

	var ShortID string
	var OriginalURL string
	var ExpireTime string
	err = stmt.QueryRow(shorturlID).Scan(&ShortID, &OriginalURL, &ExpireTime)
	defer stmt.Close()
	if err != nil {
		return false, "", "", ""
	}

	return true, ShortID, OriginalURL, ExpireTime
}

// If data expired, delete the data.
func Pg_Del(db *sql.DB, shorturlID string) error {
	stmt, err := db.Prepare("DELETE FROM shortenerdb WHERE shortid = $1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(shorturlID)
	defer stmt.Close()
	if err != nil {
		return err
	}

	return nil
}
