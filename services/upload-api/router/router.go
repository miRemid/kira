package router

import (
	"github.com/gin-gonic/gin"
	mClient "github.com/micro/go-micro/v2/client"

	"github.com/miRemid/kira/client"
	"github.com/miRemid/kira/common/middleware"
)

var (
	upload  *client.UploadClient
	fileCli *client.FileClient
)

func NewRouter(cli mClient.Client) *gin.Engine {
	upload = client.NewUploadClient(cli)
	fileCli = client.NewFileClient(cli)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	upload := router.Group("/upload", middleware.APICount("upload"), CheckFileToken(), CheckSuspend)
	{
		upload.PUT("/image", UploadFile)
	}

	return router
}
