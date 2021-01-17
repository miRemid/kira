package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("usernameValidate", usernameValidator)
		v.RegisterValidation("passwordValidate", passwordValidator)
	}

	v1 := route.Group("/user")
	{

		v1.POST("/signup", Signup)
		v1.POST("/signin", Signin)
		v1.GET("/me", GetInfo)
		v1.GET("/refreshToken", RefreshToken)
	}

	return route
}
