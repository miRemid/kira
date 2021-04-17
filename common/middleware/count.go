package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/cache/redis"
	"github.com/miRemid/kira/common"
)

func APICount(service string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conn := redis.Get()
		defer conn.Close()
		key := common.ServiceKey(service)
		conn.Do("SADD", "kira-services", key)
		conn.Do("HINCRBY", key, ctx.Request.URL.Path, 1)
		ctx.Next()
	}
}
