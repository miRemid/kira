package router

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/services/upload/config"
)

func UploadFile(ctx *gin.Context) {
	owner, _ := ctx.Get("owner")
	width := ctx.Query("width")
	height := ctx.Query("height")
	log.Println("Upload File, Owner: ", owner)
	file, meta, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Println("UploadFile, get file from form err: ", err)
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
		log.Println("UploadFile, ext error: ", fileExt)
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: "not support ext",
		})
		return
	}
	res, err := upload.UploadFile(owner.(string), fileName, fileExt, width, height, file)
	if err != nil {
		log.Println("Upload File: ", err)
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	} else if !res.Succ {
		log.Println("UploadFile: ", res.Msg)
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: res.Msg,
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code:    response.StatusOK,
		Message: res.Msg,
		Data:    res.File,
	})
}
