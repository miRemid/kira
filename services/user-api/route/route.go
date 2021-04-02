package route

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/miRemid/kira/client"
	"github.com/miRemid/kira/common/middleware"
	microClient "github.com/micro/go-micro/v2/client"
)

var (
	cli     *client.UserClient
	authCli *client.AuthClient
)

func Init(clients microClient.Client) {
	cli = client.NewUserClient(clients)
	authCli = client.NewAuthClient(clients)
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

	v1 := route.Group("/user", PrintlnPath, middleware.APICount("user"))
	{

		v1.POST("/signup", Signup)
		v1.POST("/signin", Signin)

		auth := v1.Group("/", JwtAuth(e))
		{
			auth.GET("/me", GetInfo)

			auth.DELETE("/deleteUser", DeleteUser)

			auth.GET("/getUserList", GetUserList)

			auth.POST("/updateUser", UpdateUser)
		}

	}

	return route
}
