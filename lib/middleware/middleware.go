package middleware

import (
	"net/http"
	//"net/http"
	"github.com/WeiWeiCheng123/URLshortener/model"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

var (
	db            *xorm.Engine
	IPLimitMax    int
	IPLimitPeriod int
)

func Init(database *xorm.Engine) {
	db = database
}

func Init_ip(ip_max int, ip_limit_period int) {
	IPLimitMax = ip_max
	IPLimitPeriod = ip_limit_period
}

//Limit IP usage
func IPLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := model.Redis_ip_limit(c.ClientIP(), IPLimitMax, IPLimitPeriod)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if res == -1 {
			c.String(http.StatusTooManyRequests, "Too many requests")
			return
		}
	}
}
