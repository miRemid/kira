package handler

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/cache/redis"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/services/file/config"
	"github.com/miRemid/kira/services/site/client"
	"golang.org/x/sync/errgroup"
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
	width := ctx.DefaultQuery("width", "0")
	height := ctx.DefaultQuery("height", "0")
	res, err := client.File().GetImage(fileID, width, height)
	if err != nil {
		log.Println("Get Image: ", err)
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	} else if !res.Succ {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: res.Msg,
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
		log.Println("Get API Counts: ", err)
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

func Ping(ctx *gin.Context) {
	// errorGroup
	g, _ := errgroup.WithContext(context.Background())
	var res = make(map[string]interface{})
	var lock sync.Mutex
	g.Go(func() error {
		f := client.FileCli
		resp, err := f.Ping()
		if err != nil {
			return err
		}
		lock.Lock()
		res[resp.Name] = resp.Message
		lock.Unlock()
		return nil
	})
	g.Go(func() error {
		u := client.UserCli
		resp, err := u.Ping()
		if err != nil {
			return err
		}
		lock.Lock()
		res[resp.Name] = resp.Message
		lock.Unlock()
		return nil
	})
	g.Go(func() error {
		a := client.AuthCli
		resp, err := a.Ping()
		if err != nil {
			return err
		}
		lock.Lock()
		res[resp.Name] = resp.Message
		lock.Unlock()
		return nil
	})
	g.Go(func() error {
		u := client.UploadCli
		resp, err := u.Ping()
		if err != nil {
			return err
		}
		lock.Lock()
		res[resp.Name] = resp.Message
		lock.Unlock()
		return nil
	})
	if err := g.Wait(); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusPingError,
			Error: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code: response.StatusOK,
		Data: res,
	})
}
