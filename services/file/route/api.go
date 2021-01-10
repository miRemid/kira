package route

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/services/file/config"
	"github.com/miRemid/kira/services/file/pb"
)

type Search struct {
	Offset int64 `form:"offset"`
	Limit  int64 `form:"limit"`
}

type SearchRes struct {
	Total int64      `json:"total"`
	Files []*pb.File `json:"files"`
}

func GetHistory(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	var s Search
	if err := ctx.BindQuery(&s); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: "missing params",
		})
		return
	}
	if s.Limit == 0 {
		s.Limit = 10
	}
	res, err := cli.GetHistory(token.(string), s.Limit, s.Offset)
	if err != nil || !res.Succ {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: res.Msg,
		})
		return
	}
	var sr SearchRes
	sr.Total = res.Total
	sr.Files = res.Files
	ctx.JSON(http.StatusOK, response.Response{
		Code:    response.StatusOK,
		Message: "get success",
		Data:    sr,
	})
}

func UploadFile(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	file, meta, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:  response.StatusBadParams,
			Error: "missing file",
		})
		return
	}
	defer file.Close()
	// 1. check ext
	fileName := meta.Filename
	fileExt := filepath.Ext(fileName)
	if !config.CheckExt(fileExt) {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: "not support ext",
		})
		return
	}

	res, err := cli.UploadFile(token.(string), fileName, fileExt, file)
	if err != nil || !res.Succ {
		log.Print(err)
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code:    response.StatusOK,
		Message: res.Msg,
		Data:    res.File,
	})
}

type DeleteReq struct {
	FileID string `json:"file_id" form:"file_id" binding:"required"`
}

func DeleteFile(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	var req DeleteReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: err.Error(),
		})
		return
	}
	res, err := cli.DeleteFile(token.(string), req.FileID)
	if err != nil || !res.Succ {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code:    response.StatusOK,
		Message: res.Msg,
	})
}

func GetImage(ctx *gin.Context) {
	fileID := ctx.Param("fileid")
	if fileID == "" {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: "mising fileid param",
		})
		return
	}
	res, err := cli.GetImage(fileID)
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
