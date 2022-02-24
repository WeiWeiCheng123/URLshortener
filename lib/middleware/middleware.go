package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

const (
	IPLimitPeriod     = 1800
	IPLimitTimeFormat = "2006-01-02 15:04:05"
	IPLimitMax        = 10
)

var rdb *redis.Pool

func Init(r *redis.Pool) {
	rdb = r
}

func IPLimiter(c *gin.Context) {
	fmt.Println("ip= ", c.ClientIP())
	con := rdb.Get()
	defer con.Close()

	cont, err := redis.Int(con.Do("GET", c.ClientIP()))
	fmt.Println(cont)
	if err != nil {
		con.Do("SET", c.ClientIP(), 1)
		con.Do("EXPIRE", c.ClientIP(), IPLimitPeriod)
	} else {
		con.Do("EXPIRE", c.ClientIP(), IPLimitPeriod)
		if cont > IPLimitMax {
			con.Do("INCR", c.ClientIP())
		}
	}

	if cont >= IPLimitMax {
		c.String(429, "Too many requests")
		c.Abort()
	}

	c.Next()
}
