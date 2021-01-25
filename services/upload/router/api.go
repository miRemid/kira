package router

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/services/upload/config"
)

func UploadFile(ctx *gin.Context) {
	owner, _ := ctx.Get("owner")
	file, meta, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:  response.StatusBadParams,
			Error: "missing file",
		})
		return
	}
	defer file.Close()
	fileName := meta.Filename
	fileExt := filepath.Ext(fileName)
	if !config.CheckExt(fileExt) {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: "not support ext",
		})
		return
	}
	res, err := upload.UploadFile(owner.(string), fileName, fileExt, file)
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
		Data:    res.File,
	})
}
