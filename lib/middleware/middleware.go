package middleware

import (
	"net/http"
	//"net/http"

	"github.com/WeiWeiCheng123/URLshortener/model"
	"github.com/gin-gonic/gin"
)

var (
	IPLimitMax    int
	IPLimitPeriod int
)

func Init(ip_max int, ip_limit_period int) {
	IPLimitMax = ip_max
	IPLimitPeriod = ip_limit_period
}

//Limit IP usage
func IPLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := model.Redis_ip_limit(c.ClientIP(), IPLimitMax, IPLimitPeriod)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}

		if res == -1 {
			c.String(http.StatusTooManyRequests, "Too many requests")
			c.Abort()
		}
	}
}
