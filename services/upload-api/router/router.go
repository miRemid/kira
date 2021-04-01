package router

import (
	"github.com/gin-gonic/gin"
	mClient "github.com/micro/go-micro/v2/client"

	"github.com/miRemid/kira/common/middleware"

	"github.com/miRemid/kira/client"
)

var (
	auth   *client.AuthClient
	upload *client.UploadClient
)

func NewRouter(cli mClient.Client) *gin.Engine {
	auth = client.NewAuthClient(cli)
	upload = client.NewUploadClient(cli)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	upload := router.Group("/upload", middleware.APICount("upload"))
	{
		upload.PUT("/image", UploadFile)
	}

	return router
}
