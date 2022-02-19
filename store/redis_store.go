package store

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func NewPool(addr string, max int, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     max, //Max connection
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			_, err = c.Do("AUTH", password)
			if err != nil {
				c.Close()
				return nil, err
			}
			fmt.Println("Redis connect!")
			return c, nil
		},
	}
}

func Get() redis.Conn {
	return pool.Get()
}

//Give original URL and expire time, save to Redis.
//Key: shortURL  ; 	Value: URL  ; 	TTL: expire time  ;
func Redis_Save(r *redis.Pool, shorturlID string, url string, expireTime time.Time) (string, error) {
	connections := r.Get()
	defer connections.Close()

	_, err := connections.Do("SET", shorturlID, url)
	if err != nil {
		return "", err
	}

	_, err = connections.Do("EXPIREAT", shorturlID, expireTime.Unix())
	if err != nil {
		return "", err
	}

	return shorturlID, nil
}

//Give shortURL if not expired return original URL
func Redis_Load(r *redis.Pool, shortURL string) (string, error) {
	connections := r.Get()
	defer connections.Close()

	url, err := redis.String(connections.Do("GET", shortURL))
	if err != nil {
		return "", err
	}
	return url, nil
}
