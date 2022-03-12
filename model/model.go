package model

import (
	"database/sql"
	"time"

	"github.com/gomodule/redigo/redis"
)

type ShortURL struct {
	ShortID     string
	OriginalURL string
	ExpireTime  time.Time
}
var (
	pdb *sql.DB
	rdb *redis.Pool
)  

//Connect to postgres and redis
func Init(p *sql.DB, r *redis.Pool) {
	pdb = p
	rdb = r
}