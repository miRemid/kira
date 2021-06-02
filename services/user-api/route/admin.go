package route

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/proto/pb"
)

func GetAnonyList(ctx *gin.Context) {
	var req pb.GetAnonyFilesReq
	if err := ctx.BindQuery(&req); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: err.Error(),
		})
		return
	}
	log.Println(req.Limit, req.Offset)
	resp, err := fileCli.Service.GetAnonyFiles(ctx, &req)
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
	if err := ctx.BindQuery(&req); err != nil {
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
