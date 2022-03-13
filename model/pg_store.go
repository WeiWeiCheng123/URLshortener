package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/WeiWeiCheng123/URLshortener/lib/function"
	_ "github.com/lib/pq"
)

//Connect to postgres
func Connect_Pg(connect string) *sql.DB {
	db, err := sql.Open("postgres", connect)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db
}

//Give original URL and expire time, save to Postgres.
func Pg_Save(url string, expireTime time.Time) (string, error) {
	stmt, err := pdb.Prepare("INSERT INTO shortener(shortid,originalurl,expiretime) VALUES($1,$2,$3)")
	defer stmt.Close()

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	shortID := function.Generator()
	_, err = stmt.Exec(shortID, url, expireTime)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return shortID, nil
}

// If there is not exist, return false, otherwise return true
func Pg_Load(shorturlID string) (bool, string, time.Time) {
	stmt, err := pdb.Prepare("SELECT originalurl, expiretime FROM shortener WHERE shortid = $1")
	defer stmt.Close()

	if err != nil {
		return false, "", time.Time{}
	}

	data := ShortURL{}
	err = stmt.QueryRow(shorturlID).Scan(&data.OriginalURL, &data.ExpireTime)
	if err != nil {
		return false, "", time.Time{}
	}

	return true, data.OriginalURL, data.ExpireTime
}

// If data expired, delete the data.
func Pg_Del(shorturlID string) error {
	stmt, err := pdb.Prepare("DELETE FROM shortener WHERE shortid = $1")
	defer stmt.Close()

	if err != nil {
		return err
	}

	_, err = stmt.Exec(shorturlID)
	if err != nil {
		return err
	}

	return nil
}

func Pg_Del_Exp() {
	fmt.Println("Cron Job start", time.Now())
	stmt, _ := pdb.Prepare("DELETE FROM shortener WHERE expireTime < $1")
	defer stmt.Close()

	stmt.Exec(time.Now())
	fmt.Println("Cron Job done")
}
