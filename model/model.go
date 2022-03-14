package model

import (
	"time"
)

type ShortURL struct {
	ShortID     string `xorm:pk json:shortid`
	OriginalURL string `json:original_url`
	ExpireTime  time.Time `json:expire_time`
}

func (u *ShortURL) TableName() string {
	return "ShortURL"
}