package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/services/site/handler"
)

func New() *gin.Engine {
	route := gin.New()
	route.Use(gin.Logger())
	route.Use(gin.Recovery())

	v1 := route.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "ping")
		})
	}

	route.GET("/image", handler.GetImage)

	return route
}
