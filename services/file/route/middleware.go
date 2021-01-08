package route

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/miRemid/kira/common/response"
)

func CheckToken(ctx *gin.Context) {
	if token := ctx.Query("token"); token == "" {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Code:  response.StatusUnauthorized,
			Error: "missing token params",
		})
		ctx.Abort()
		return
	} else {
		ctx.Set("token", token)
	}
	ctx.Next()
}
