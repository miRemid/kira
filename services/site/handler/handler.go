package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/cache/redis"
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

func GetAPICounts(ctx *gin.Context) {
	// 1. Get services list
	res, err := redis.Get().Do("SMEMBERS", "kira")
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}
	data := make(map[string]map[string]string)
	items, _ := res.([]interface{})
	for _, item := range items {
		path, _ := item.([]byte)
		key := string(path)
		arr := strings.Split(key, "-")
		service, router := arr[0], arr[1]
		if _, ok := data[service]; !ok {
			data[service] = make(map[string]string)
		}
		count, _ := redis.Get().Do("GET", key)
		data[service][router] = string(count.([]byte))
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code: response.StatusOK,
		Data: data,
	})
}
