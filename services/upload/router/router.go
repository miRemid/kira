package router

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"

	authClient "github.com/miRemid/kira/services/auth/client"
	uploadClient "github.com/miRemid/kira/services/upload/client"
)

var (
	auth   *authClient.AuthClient
	upload *uploadClient.UploadClient
)

func NewRouter(cli client.Client) *gin.Engine {
	auth = authClient.NewAuthClient(cli)
	upload = uploadClient.NewUploadClient(cli)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("/v1")
	{
		// v1/user/upload
		user := v1.Group("/", CheckToken)
		{
			user.PUT("/upload", UploadFile)
		}
	}

	return router
}
