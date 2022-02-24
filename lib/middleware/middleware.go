package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

const (
	IPLimitPeriod     = 1800
	IPLimitTimeFormat = "2006-01-02 15:04:05"
	IPLimitMax        = 500
)
var rdb *redis.Pool

func Init(r *redis.Pool) {
	rdb = r
}

func IPLimiter(c *gin.Context) {
	fmt.Println("ip= ", c.ClientIP())
	c.Next()
}

