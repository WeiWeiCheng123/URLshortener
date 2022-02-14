package store

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/WeiWeiCheng123/URLshortener/function"
	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func NewPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10, //Max connection
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			_, err = c.Do("AUTH", "password")
			if err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}

func Get() redis.Conn {
	return pool.Get()
}

//Check the id is exist or not
func CheckId(r *redis.Pool, id uint64) bool {
	connections := r.Get()
	defer connections.Close()

	is_exists, _ := redis.String(connections.Do("GET", strconv.FormatUint(id, 10)))
	if is_exists == "" {
		return false
	}

	return true
}

//Give original URL and expire time, save to Redis
func Save(r *redis.Pool, url string, expireTime time.Time) (string, error) {
	connections := r.Get()
	defer connections.Close()

	var id uint64

	for exist := true; exist; exist = CheckId(r, id) {
		id = rand.Uint64()
	}

	_, err := connections.Do("SET", strconv.FormatUint(id, 10), url)
	if err != nil {
		fmt.Println("set:", err)
		return "", err
	}

	_, err = connections.Do("EXPIREAT", strconv.FormatUint(id, 10), expireTime.Unix())
	if err != nil {
		fmt.Println("exp", err)
		return "", err
	}

	return function.Encode(id), nil
}

//Give shortURL return original URL
func Load(r *redis.Pool, shortURL string) (string, error) {
	id, err := function.Decode(shortURL)
	connections := r.Get()
	if err != nil {
		return "", err
	}

	url, err := redis.String(connections.Do("GET", strconv.FormatUint(id, 10)))
	if err != nil {
		return "", err
	}
	return url, nil
}