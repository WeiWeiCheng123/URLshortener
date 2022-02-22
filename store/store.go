package store

import (
	"database/sql"

	"github.com/gomodule/redigo/redis"
)

var pdb *sql.DB
var rdb *redis.Pool

type ShortURL struct {
	ShortID     string
	OriginalURL string
	ExpireTime  string
}

//Connect to postgres and redis
func Init(postgres_db *sql.DB, redis_db *redis.Pool) {
	pdb = postgres_db
	rdb = redis_db
}
