package route

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/miRemid/kira/common/response"
)

func CheckToken(ctx *gin.Context) {
	owner := ""
	if token := ctx.Query("token"); token == "" {
		owner = ctx.ClientIP()
	} else {
		res, err := auth.FileToken(token)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.Response{
				Code:  response.StatusInternalError,
				Error: err.Error(),
			})
			return
		}
		owner = res.UserID
	}
	ctx.Set("owner", owner)
	ctx.Next()
}
