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
	return "Shortener"
}
