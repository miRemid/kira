package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/cache/redis"
)

func APICount(service string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conn := redis.Get()
		defer conn.Close()
		conn.Do("SADD", "kira-services", service)
		conn.Do("HINCRBY", service, ctx.Request.URL.Path, 1)
		ctx.Next()
	}
}
