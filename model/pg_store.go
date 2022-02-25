package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/WeiWeiCheng123/URLshortener/pkg/function"
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
	stmt, err := pdb.Prepare("INSERT INTO shortenerdb(shortid,originalurl,expiretime) VALUES($1,$2,$3)")
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	shortID := function.Id()
	_, err = stmt.Exec(shortID, url, expireTime)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return shortID, nil
}

// If there is not exist, return false, otherwise return true
func Pg_Load(shorturlID string) (bool, string, string, time.Time) {
	stmt, err := pdb.Prepare("SELECT shortid, originalurl, expiretime FROM shortenerdb WHERE shortid = $1")
	if err != nil {
		return false, "", "", time.Time{}
	}
	data := ShortURL{}
	err = stmt.QueryRow(shorturlID).Scan(&data.ShortID, &data.OriginalURL, &data.ExpireTime)
	defer stmt.Close()
	if err != nil {
		return false, "", "", time.Time{}
	}

	return true, data.ShortID, data.OriginalURL, data.ExpireTime
}

// If data expired, delete the data.
func Pg_Del(shorturlID string) error {
	stmt, err := pdb.Prepare("DELETE FROM shortenerdb WHERE shortid = $1")
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

func Pg_Del_Exp() {
	fmt.Println("Cron Job start")
	stmt, _ := pdb.Prepare("DELETE FROM shortenerdb WHERE expireTime < $1")
	stmt.Exec(time.Now())
	stmt.Close()
	fmt.Println("Cron Job done")
}
