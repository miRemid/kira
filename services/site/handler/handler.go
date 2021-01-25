package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/services/file/config"
	"github.com/miRemid/kira/services/site/client"
)

func GetImage(ctx *gin.Context) {
	fileID := ctx.Query("id")
	if fileID == "" {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: "mising fileid param",
		})
		return
	}

	res, err := client.File().GetImage(fileID)
	if err != nil || !res.Succ {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}
	// 写文件
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Header().Add("Content-Type", config.ContentType(res.FileExt))
	ctx.Writer.Write(res.Image)
}
