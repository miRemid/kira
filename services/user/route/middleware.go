package route

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/services/auth/pb"
)

func parseToken(header string) (res *pb.ValidResponse, err error) {
	split := strings.Split(header, " ")
	if len(split) != 2 {
		return nil, errors.New("invalid token struct")
	}
	if split[0] != "Bearer" {
		return nil, errors.New("invalid prefix")
	}
	log.Println(split[1])
	return authCli.Valid(split[1])
}

func JwtAuth(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		if path == "/v1/signup" || path == "/v1/signin" {
			ctx.Next()
			return
		}

		header := ctx.Request.Header.Get("Authorization")
		if header == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		res, err := parseToken(header)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if res.Expired {
			ctx.JSON(http.StatusOK, response.Response{
				Code:  response.StatusExpired,
				Error: "signin expired",
			})
			return
		}

		if allow, err := enforcer.Enforce(res.UserRole, path, method); err != nil {
			ctx.AbortWithError(http.StatusForbidden, err)
			return
		} else if !allow {
			ctx.AbortWithError(http.StatusForbidden, errors.New("no permission"))
			return
		}

		ctx.Request.Header.Set("userid", res.UserID)
		ctx.Request.Header.Set("userRole", res.UserRole)
		ctx.Next()
	}
}
