package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	if gin.Mode() == gin.ReleaseMode {
		config.AllowOrigins = []string{"*"}
	} else {
		config.AllowOrigins = []string{"*"}
	}
	config.AllowCredentials = true
	return cors.New(config)
}
