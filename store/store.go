package store

import (
	"database/sql"
	"time"

	"github.com/gomodule/redigo/redis"
)

func Save(r *redis.Pool, p *sql.DB, url string, expireTime time.Time) {
	redis_con := r.Get()
	defer redis_con.Close()
	
}
