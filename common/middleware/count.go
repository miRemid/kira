package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/cache/redis"
)

func APICount(service string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		service := fmt.Sprintf("%s-%s", service, ctx.Request.URL.Path)
		redis.Get().Do("SADD", "kira", service)
		redis.Get().Do("incr", service)
		ctx.Next()
	}
}
