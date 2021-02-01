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

	v1 := router.Group("/v1", middleware.APICount("upload"))
	{
		// v1/user/upload
		user := v1.Group("/", CheckToken)
		{
			user.PUT("/upload", UploadFile)
		}
	}

	return router
}
