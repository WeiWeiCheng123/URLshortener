package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

var rdb *redis.Pool
var IPLimitMax int
var IPLimitPeriod int

func Init(r *redis.Pool, ip_max int, ip_limit_period int) {
	rdb = r
	IPLimitMax = ip_max
	IPLimitPeriod = ip_limit_period
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
		if cont < IPLimitMax {
			con.Do("INCR", c.ClientIP())
		}
	}

	if cont >= IPLimitMax {
		c.String(429, "Too many requests")
		c.Abort()
	}

	c.Next()
}
