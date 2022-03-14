package model

import (
	"time"
)

type Shortener struct {
	ShortId     string    `xorm:pk json:short_id`
	OriginalUrl string    `json:original_url`
	ExpireTime  time.Time `json:expire_time`
}

func (u *Shortener) TableName() string {
	return "shortener"
}


//docker rmi url-shortener
//docker-compose down
//docker-compose up
//curl -X POST -H "Content-Type:application/json" -d '{"url":"https://www.dcard.tw/f","expireAt":"2023-01-01T09:00:41Z"}' http://localhost:8080/api/v1/urls
//curl -L -X GET "http://localhost:8080/WmBfwUK"

/*
git add .
git commit -m "fix bug"
git push origin v2
*/