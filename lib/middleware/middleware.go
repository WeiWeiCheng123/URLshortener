package middleware

import (
	"fmt"
	"net/http"

	"github.com/WeiWeiCheng123/URLshortener/lib/lua"
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

func IPLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		con := rdb.Get()
		defer con.Close()

		script := redis.NewScript(1, lua.IP_script)
		res, err := redis.Int(script.Do(con, c.ClientIP(), IPLimitMax, IPLimitPeriod))
		/*
			cont, err := redis.Int(con.Do("GET", c.ClientIP()))
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
				c.String(http.StatusTooManyRequests, "Too many requests")
				c.Abort()
			}
		*/
		if err != redis.ErrNil {
			fmt.Println("err: ", res, " ; ", err.Error())
		}
		if res == -1 {
			c.String(http.StatusTooManyRequests, "Too many requests")
			c.Abort()
		}
	}
}
