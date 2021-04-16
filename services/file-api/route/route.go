package route

import (
	"github.com/casbin/casbin/v2"
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

func Route(e *casbin.Enforcer) *gin.Engine {

	route := gin.New()
	route.Use(gin.Logger())
	route.Use(gin.Recovery())

	file := route.Group("/file", middleware.APICount("file"))
	{
		file.GET("/getRandomFiles", GetRandomFile)
		file.GET("/getHotLikeRank", GetHotLikeRank)
		token := file.Group("/", CheckToken)
		{
			token.GET("/history", GetHistory)
			token.DELETE("/delete", DeleteFile)
			token.GET("/detail", GetDetail)
			token.GET("/refreshToken", RefreshToken)
		}

		auth := file.Group("/", middleware.JwtAuth(auth, e))
		{
			auth.GET("/getToken", GetUserToken)
			auth.POST("/like", LikeOrDislike)
		}
	}

	return route
}
