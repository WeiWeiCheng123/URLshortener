package middleware

import (
	"net/http"

	"github.com/WeiWeiCheng123/URLshortener/lib/constant"
	"github.com/WeiWeiCheng123/URLshortener/lib/function"
	"github.com/WeiWeiCheng123/URLshortener/lib/lua"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/gomodule/redigo/redis"
)

var (
	rdb           *redis.Pool
	db            *xorm.Engine
	IPLimitMax    int
	IPLimitPeriod int
)

func Init(database *xorm.Engine, r *redis.Pool) {
	db = database
	rdb = r
}

func Init_ip(ip_max int, ip_limit_period int) {
	IPLimitMax = ip_max
	IPLimitPeriod = ip_limit_period
}

// do nothing and provide an injection of database object only
func Plain() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constant.DB, db)
		c.Set(constant.StatusCode, nil)
		c.Set(constant.Error, nil)
		c.Set(constant.Output, nil)
		c.Next()

		statusCode := c.GetInt(constant.StatusCode)
		err := c.MustGet(constant.Error)
		output := c.MustGet(constant.Output)
		if err != nil {
			c.String(statusCode, err.(error).Error())
		} else {
			c.JSON(statusCode, output)
		}
	}
}

// check the shortid and provide an injection of database and redis objects
func TX() gin.HandlerFunc {
	return func(c *gin.Context) {
		shortID := c.Param("shortID")
		if len(shortID) != 7 {
			c.String(http.StatusNotFound, "this short URL is not existed or expired")
			c.Abort()
			return
		}
		if err := function.ShortID_legal(shortID); err != nil {
			c.String(http.StatusNotFound, "this short URL is not existed or expired")
			c.Abort()
			return
		}

		c.Set(constant.ShortID, c.Param("shortID"))
		c.Set(constant.DB, db)
		c.Set(constant.Cache, rdb)
		c.Set(constant.StatusCode, nil)
		c.Set(constant.Error, nil)
		c.Set(constant.Output, nil)
		c.Next()

		statusCode := c.GetInt(constant.StatusCode)
		output := c.MustGet(constant.Output)
		err := c.MustGet(constant.Error)
		if err != nil {
			c.String(statusCode, err.(error).Error())
		} else {
			c.Redirect(statusCode, output.(string))
		}
	}
}

// limit ip usage
func IPLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		connections := rdb.Get()
		defer connections.Close()

		script := redis.NewScript(1, lua.IP_script)
		res, err := redis.Int(script.Do(connections, c.ClientIP(), IPLimitMax, IPLimitPeriod))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}
		// exceed the limit
		if res == -1 {
			c.String(http.StatusTooManyRequests, "too many requests")
			c.Abort()
			return
		}
	}
}