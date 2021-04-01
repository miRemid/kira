package route

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/client"
	"github.com/miRemid/kira/common/middleware"
	microClient "github.com/micro/go-micro/v2/client"
)

var (
	cli  *client.FileClient
	auth *client.AuthClient
)

func Init(clients microClient.Client) {
	cli = client.NewFileClient(clients)
	auth = client.NewAuthClient(clients)
}

func Route() *gin.Engine {

	route := gin.New()
	route.Use(gin.Logger())
	route.Use(func(ctx *gin.Context) {
		log.Println(ctx.Request.RemoteAddr, ctx.Request.URL.Path)
	})
	route.Use(gin.Recovery())

	route.Use(middleware.CORS())

	file := route.Group("/file", CheckToken, middleware.APICount("file"))
	{
		file.GET("/history", GetHistory)
		file.DELETE("/delete", DeleteFile)
		file.GET("/detail", GetDetail)
	}

	return route
}
