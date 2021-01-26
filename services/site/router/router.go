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

	v1 := route.Group("/v1/site", middleware.APICount("site"))
	{
		v1.GET("/info", handler.GetAPICounts)
	}

	route.GET("/image", handler.GetImage)

	return route
}