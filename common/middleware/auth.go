package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/client"
	"github.com/miRemid/kira/common/response"
)

func parseToken(header string) (string, error) {
	split := strings.Split(header, " ")
	if len(split) != 2 {
		return "", errors.New("invalid token struct")
	}
	if split[0] != "Bearer" {
		return "", errors.New("invalid prefix")
	}
	return split[1], nil
}

func JwtAuth(authCli *client.AuthClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.Request.Header.Get("Authorization")
		if header == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := parseToken(header)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		res, err := authCli.Valid(token)
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

		ctx.Request.Header.Set("userid", res.UserID)
		ctx.Request.Header.Set("userRole", res.UserRole)
		ctx.Next()
	}
}
