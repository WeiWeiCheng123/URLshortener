package model

import (
	"time"
)

type Shortener struct {
	ShortID     string    `json:"shortid" xorm:"pk"`
	OriginalURL string    `json:"url"`
	ExpireTime  time.Time `json:"expireAt"`
}

func (s Shortener) TableName() string {
	return "shortener"
}
