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

func GetHotLikeRank(ctx *gin.Context) {
	res, err := cli.Service.GetHotLikeRank(ctx, &pb.Empty{})
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
			"files": res.Files,
		},
	})
}
