package handler

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/miRemid/kira/cache/redis"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/file/config"
	"github.com/miRemid/kira/services/site/client"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/sync/errgroup"
)

func GetImage(ctx *gin.Context) {
	var in = new(pb.GetImageReq)
	if err := ctx.ShouldBind(in); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: "mising params",
		})
		return
	}
	res, err := client.File().Service.GetImage(ctx, in)
	if err != nil {
		log.Println("Get Image: ", err)
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	} else if !res.Succ {
		ctx.JSON(http.StatusNotFound, response.Response{
			Code:  response.StatusInternalError,
			Error: res.Msg,
		})
		return
	}
	ctx.Writer.Header().Set("Content-Type", config.ContentType(res.FileExt))
	ctx.Writer.Header().Set("Content-Disposition", "filename="+res.FileName)
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(res.Image)
}

func DownloadImage(ctx *gin.Context) {
	var in = new(pb.GetImageReq)
	if err := ctx.ShouldBind(in); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: "mising params",
		})
		return
	}
	res, err := client.File().Service.GetImage(ctx, in)
	if err != nil {
		log.Println("Get Image: ", err)
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	} else if !res.Succ {
		ctx.JSON(http.StatusNotFound, response.Response{
			Code:  response.StatusInternalError,
			Error: res.Msg,
		})
		return
	}
	ctx.Writer.Header().Set("Content-Type", config.ContentType(res.FileExt))
	ctx.Writer.Header().Set("Content-Disposition", "attachment;filename="+res.FileName)
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(res.Image)
}

func GetAPICounts(ctx *gin.Context) {
	// 1. Get services list
	conn := redis.Get()
	defer conn.Close()
	services, err := redigo.Strings(conn.Do("SMEMBERS", "kira-services"))
	if err != nil {
		log.Println("Get API Counts: ", err)
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}
	// 2. Get each service's api calls count
	data := make(map[string]interface{})
	sum := make(map[string]int64)
	for _, key := range services {
		// key is the service's name, eg: file
		// 3. get file service's all count
		// all the counts store in a hashmap named file(key)
		strMap, err := redigo.Int64Map(conn.Do("HGETALL", key))
		if err != nil {
			log.Printf("Get %s counts err: ", err)
		} else {
			s := int64(0)
			for _, v := range strMap {
				s += v
			}
			sum[key] = s
			data[key] = strMap
		}
	}
	data["sum"] = sum
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

func GetHostStat(ctx *gin.Context) {
	mstat, _ := mem.VirtualMemory()
	sstat, _ := mem.SwapMemory()
	dstat, _ := disk.Usage("/")
	cstat, _ := cpu.Percent(time.Second, true)
	cpuStat := 0.0
	for _, v := range cstat {
		cpuStat += v
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: response.StatusOK,
		Data: gin.H{
			"Memory": mstat.UsedPercent,
			"Disk":   dstat.UsedPercent,
			"Cpu":    cpuStat,
			"Swap":   sstat.UsedPercent,
		},
	})
}
