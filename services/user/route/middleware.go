package route

import "github.com/gin-gonic/gin"

func GetUserID(ctx *gin.Context) {
	// user id will be set in the request header
	// X-XXID
	userid := ctx.GetHeader("X-XXID")
	ctx.Set("userid", userid)
	ctx.Next()
}
