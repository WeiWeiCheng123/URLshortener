package model

import (
	"time"

	"github.com/WeiWeiCheng123/URLshortener/lib/lua"
	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

//Connect to redis
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
			return c, nil
		},
	}
}

func Get() redis.Conn {
	return pool.Get()
}

//Give original URL and expire time, save to Redis.
//Key: shortURL  ; 	Value: URL  ; 	TTL: expire time  ;
func Redis_Save(shortURL string, url string, expireTime time.Time) error {
	connections := rdb.Get()
	defer connections.Close()
	
	_, err := connections.Do("SETEX", shortURL, int(expireTime.Sub(time.Now()).Seconds()), url)
	if err != nil {
		return err
	}

	return nil
}

//Give shortURL if not expired return original URL
func Redis_Load(shortURL string) (string, error) {
	connections := rdb.Get()
	defer connections.Close()

	url, err := redis.String(connections.Do("GET", shortURL))
	if err != nil {
		return "", err
	}

	return url, nil
}

//Give the not existing shortURL and set it into Redis with value NotExist.
//To prevent too many users trying to access with a non-existent shorten URL.
func Redis_Set_NotExist(shortURL string) error {
	connections := rdb.Get()
	defer connections.Close()

	_, err := connections.Do("SETEX", shortURL, 300, "NotExist")
	if err != nil {
		return err
	}

	return nil
}

//Give IP, limit max and limit period and use lua set it into Redis.
//Limit IP usage
func Redis_ip_limit(ip string, max int, period int) (int, error) {
	connections := rdb.Get()
	defer connections.Close()

	script := redis.NewScript(1, lua.IP_script)
	res, err := redis.Int(script.Do(connections, ip, max, period))
	return res, err
}
