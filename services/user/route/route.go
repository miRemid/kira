package route

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/miRemid/kira/common/middleware"
	authClient "github.com/miRemid/kira/services/auth/client"
	"github.com/miRemid/kira/services/user/client"
	microClient "github.com/micro/go-micro/v2/client"
)

var (
	cli     *client.UserClient
	authCli *authClient.AuthClient
)

func Init(clients microClient.Client) {
	cli = client.NewUserClient(clients)
	authCli = authClient.NewAuthClient(clients)
}

func Route(e *casbin.Enforcer) *gin.Engine {
	route := gin.New()
	route.Use(gin.Logger())
	route.Use(gin.Recovery())

	route.Use(middleware.CORS())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("usernameValidate", usernameValidator)
		v.RegisterValidation("passwordValidate", passwordValidator)
	}

	v1 := route.Group("/v1/user")
	{

		v1.POST("/signup", Signup)
		v1.POST("/signin", Signin)

		auth := v1.Group("/", JwtAuth(e))
		{
			auth.GET("/me", GetInfo)
		}

	}

	return route
}
