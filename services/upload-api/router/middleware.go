package router

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/response"
)

func CheckToken(ctx *gin.Context) {
	owner := ""
	if token := ctx.Query("token"); token == "" {
		owner = ctx.ClientIP()
		ctx.Set(common.AnonymousKey, true)
	} else {
		userid, err := auth.FileToken(token)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.Response{
				Code:  response.StatusInternalError,
				Error: err.Error(),
			})
			return
		}
		owner = userid.UserID
		ctx.Set(common.AnonymousKey, false)
	}
	log.Println("Owner: ", owner)
	ctx.Set("owner", owner)
	ctx.Next()
}
