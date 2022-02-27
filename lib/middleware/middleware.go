package middleware

import (
	"net/http"
	//"net/http"

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

//Limit IP usage
func IPLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		connections := rdb.Get()
		defer connections.Close()

		script := redis.NewScript(1, lua.IP_script)
		res, err := redis.Int(script.Do(connections, c.ClientIP(), IPLimitMax, IPLimitPeriod))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		if res == -1 {
			c.String(http.StatusTooManyRequests, "Too many requests")
			c.Abort()
		}
	}
}
