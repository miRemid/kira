package route

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common/middleware"
	"github.com/miRemid/kira/services/file/client"
	microClient "github.com/micro/go-micro/v2/client"
)

var (
	cli *client.FileClient
)

func Init(clients microClient.Client) {
	cli = client.NewFileClient(clients)
}

func Route() *gin.Engine {

	route := gin.New()
	route.Use(gin.Logger())
	route.Use(func(ctx *gin.Context) {
		log.Println(ctx.Request.RemoteAddr, ctx.Request.URL.Path)
	})
	route.Use(gin.Recovery())

	route.Use(middleware.CORS())

	route.GET("/file/image/:fileid", GetImage)

	file := route.Group("/v1/file", CheckToken)
	{
		file.GET("/history", GetHistory)
		file.PUT("/upload", UploadFile)
		file.DELETE("/delete", DeleteFile)
		file.GET("/detail", GetDetail)
	}

	return route
}
