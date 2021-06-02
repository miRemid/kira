package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/proto/pb"
)

func GetAnonyList(ctx *gin.Context) {
	resp, err := fileCli.Service.GetAnonyFiles(ctx, &pb.GetAnonyFilesReq{})
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
			"total": resp.Total,
			"files": resp.Files,
		},
	})
}

func DeleteAnony(ctx *gin.Context) {
	var req pb.DeleteAnonyReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code: response.StatusBadParams,
		})
		return
	}
	_, err := fileCli.Service.DeleteAnonyFile(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code: response.StatusOK,
	})
}
