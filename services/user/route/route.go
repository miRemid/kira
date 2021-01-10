package route

import (
	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/services/user/client"
	microClient "github.com/micro/go-micro/v2/client"
)

var (
	cli *client.UserClient
)

func Init(clients microClient.Client) {
	cli = client.NewUserClient(clients)
}

func Route() *gin.Engine {
	route := gin.New()
	route.Use(gin.Logger())
	route.Use(gin.Recovery())

	v1 := route.Group("/api/v1")
	{

		v1.POST("/signup", Signup)
		v1.POST("/signin", Signin)

		user := v1.Group("/user")
		user.GET("/me", GetInfo)
	}

	return route
}
