package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common/response"
)

func Casbin(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		role := ctx.GetHeader("userRole")

		log.Println(path, method, role)

		if path == "/v1/user/signup" || path == "/v1/user/signin" {
			ctx.Next()
			return
		}

		if allow, err := enforcer.Enforce(role, path, method); err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, response.Response{
				Code:  response.StatusForbidden,
				Error: err.Error(),
			})
			return
		} else if !allow {
			ctx.AbortWithStatusJSON(http.StatusForbidden, response.Response{
				Code:  response.StatusForbidden,
				Error: errors.New("no permission").Error(),
			})
			return
		}
		ctx.Next()
	}
}
