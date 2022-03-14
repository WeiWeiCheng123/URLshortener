package model

import "time"

type Shortener struct {
	ShortID     string    `json:"shortid" xorm:"pk"`
	OriginalURL string    `json:"original_url"`
	ExpireTime  time.Time `json:"expire_time"`
}

func (s Shortener) TableName() string {
	return "shortener"
}
