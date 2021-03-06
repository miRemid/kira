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
	fileCli *client.FileClient
)

func Init(clients microClient.Client) {
	cli = client.NewUserClient(clients)
	authCli = client.NewAuthClient(clients)
	fileCli = client.NewFileClient(clients)
}

func Route(e *casbin.Enforcer) *gin.Engine {
	route := gin.New()
	route.Use(gin.Logger())
	route.Use(gin.Recovery())

	// route.Use(middleware.CORS())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("usernameValidate", usernameValidator)
		v.RegisterValidation("passwordValidate", passwordValidator)
	}

	v1 := route.Group("/user", middleware.APICount("user"))
	{

		v1.POST("/signup", Signup)
		v1.POST("/signin", Signin)
		v1.GET("/userInfo/:userName", GetUserInfoFromRedis, GetInfo)
		v1.POST("/forgetPassword", ForgetPassword)
		v1.POST("/modifyPassword", ModifyPassword)

		auth := v1.Group("/", middleware.JwtAuth(authCli), middleware.Casbin(e))
		{
			auth.POST("/changePassword", ChangePassword)
			auth.GET("/me", GetMe)
			auth.DELETE("/deleteAccount", DeleteAccount)
			auth.POST("/bindMail", BindMail)
			auth.GET("/vertifyMail", VertifyMail)

			admin := auth.Group("/admin")
			{
				admin.DELETE("/deleteUserFile", DeleteUserFile)
				admin.GET("/getUserList", GetUserList)
				admin.POST("/updateUserStatus", UpdateUser)

				admin.GET("/getAnonyList", GetAnonyList)
				admin.DELETE("/deleteAnony", DeleteAnony)
			}
		}
	}
	return route
}
