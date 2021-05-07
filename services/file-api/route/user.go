package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/proto/pb"
)

func GetUserToken(ctx *gin.Context) {
	userid := ctx.GetHeader("userid")

	res, err := cli.GetToken(userid)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code: response.StatusOK,
		Data: gin.H{
			"token": res.Token,
		},
	})
}

func LikeOrDislike(ctx *gin.Context) {
	userid := ctx.GetHeader("userid")
	var req = new(pb.FileLikeReq)
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: err.Error(),
		})
		return
	}
	req.Userid = userid
	res, err := cli.Service.LikeOrDislike(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code:  response.StatusCode(res.Code),
		Error: res.Message,
	})
}

func GetLikes(ctx *gin.Context) {
	userid := ctx.GetHeader("userid")
	var req = new(pb.GetLikesReq)
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: err.Error(),
		})
		return
	}
	req.Userid = userid
	if req.Limit == 0 {
		req.Limit = 5
	}
	res, err := cli.Service.GetLikes(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code: response.StatusOK,
		Data: gin.H{
			"Total": res.Total,
			"Files": res.Files,
		},
	})
}
