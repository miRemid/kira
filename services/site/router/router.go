package router

import (
	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common/middleware"
	"github.com/miRemid/kira/services/site/handler"
)

func New() *gin.Engine {
	route := gin.New()
	route.Use(gin.Logger())
	route.Use(gin.Recovery())
	// route.Use(middleware.CORS())

	v1 := route.Group("/site", middleware.APICount("site"))
	{
		v1.GET("/info", handler.GetAPICounts)
		v1.GET("/ping", handler.Ping)
	}

	route.GET("/image", handler.GetImage)

	return route
}
