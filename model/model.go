package model

import (
	"time"
)

type Shortener struct {
	ShortId     uint64       `gorm:"pk" json:"short_id"`
	OriginalUrl string    `gorm:"type:varchar(500)" json:"original_url"`
	ExpireTime  time.Time `json:"expire_time"`
}

func (u *Shortener) TableName() string {
	return "shortener"
}
