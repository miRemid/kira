package route

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/client"
	"github.com/miRemid/kira/common/middleware"
	md "github.com/miRemid/kira/services/file-api/middleware"
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

func Route(e *casbin.Enforcer) *gin.Engine {

	route := gin.New()
	route.Use(gin.Logger())
	route.Use(gin.Recovery())

	file := route.Group("/file", middleware.APICount("file"))
	{
		normal := file.Group("/", md.CheckFileToken(true))
		{
			normal.GET("/getRandomFiles", GetRandomFile)
			normal.GET("/getUserImages/:userName", GetUserImages)
		}

		token := file.Group("/", md.CheckFileToken(false))
		{
			token.GET("/history", GetHistory)
			token.DELETE("/delete", DeleteFile)
			token.GET("/detail", GetDetail)
			token.GET("/refreshToken", RefreshToken)
		}

		auth := file.Group("/", md.JwtAuth(auth), middleware.Casbin(e))
		{
			auth.GET("/getToken", GetUserToken)
			auth.POST("/like", LikeOrDislike)
			auth.GET("/getLikes", GetLikes)
		}
	}

	return route
}
