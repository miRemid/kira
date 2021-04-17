package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/response"
)

func CheckFileToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string = ""
		if t := ctx.Query("token"); t != "" {
			token = t
		} else {
			if tt := ctx.GetHeader(common.FileTokenHeader); tt != "" {
				token = tt
			}
		}
		if token == "" {
			token = common.AnonyToken
		}
		ctx.Request.Header.Set(common.FileTokenHeader, token)
		ctx.Next()
	}
}

func CheckSuspend(ctx *gin.Context) {
	if token := ctx.GetHeader(common.FileTokenHeader); token != common.AnonyToken {
		res, err := fileCli.CheckStatus(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, response.Response{
				Code:  response.StatusInternalError,
				Error: err.Error(),
			})
			return
		}
		if res.Status != 1 {
			ctx.AbortWithStatusJSON(http.StatusOK, response.Response{
				Code:  response.StatusUserSuspend,
				Error: "user upload function suspend",
			})
			return
		}
	}
	ctx.Next()
}
