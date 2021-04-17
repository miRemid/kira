package route

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/proto/pb"
)

type Search struct {
	Offset int64 `json:"offset" form:"offset"`
	Limit  int64 `json:"limit" form:"limit"`
}

func GetHistory(ctx *gin.Context) {
	var req = new(pb.GetHistoryReq)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: "missing params",
		})
		return
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	token := ctx.GetHeader(common.FileTokenHeader)
	req.Token = token
	res, err := cli.Service.GetHistory(ctx, req)
	if err != nil {
		log.Println("Get Histroy: ", err)
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}
	if !res.Succ {
		log.Println("Get Histroy: ", res.Msg)
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: res.Msg,
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code:    response.StatusOK,
		Message: "get success",
		Data: gin.H{
			"total": res.Total,
			"files": res.Files,
		},
	})
}

type DeleteReq struct {
	FileID string `json:"file_id" form:"file_id" binding:"required"`
}

func DeleteFile(ctx *gin.Context) {
	var req DeleteReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: err.Error(),
		})
		return
	}
	token := ctx.GetHeader(common.FileTokenHeader)
	res, err := cli.DeleteFile(token, req.FileID)
	if err != nil {
		log.Println("Delete File: ", err)
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}
	if !res.Succ {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: res.Msg,
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code:    response.StatusOK,
		Message: res.Msg,
	})
}

type GetDetailReq struct {
	FileID string `json:"file_id" form:"file_id" binding:"required"`
}

func GetDetail(ctx *gin.Context) {
	var req GetDetailReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: err.Error(),
		})
		return
	}
	res, err := cli.GetDetail(req.FileID)
	if err != nil {
		log.Println("Get Detail: ", err)
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
	} else {
		ctx.JSON(http.StatusOK, response.Response{
			Code: response.StatusOK,
			Data: res.File,
		})
	}
}

func RefreshToken(ctx *gin.Context) {
	token := ctx.GetHeader(common.FileTokenHeader)
	if token == "" {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusNeedToken,
			Error: "need token",
		})
		return
	}
	res, err := cli.RefreshToken(token)
	if err != nil {
		log.Println("Refresh Token Err: ", err)
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
	ctx.JSON(http.StatusOK, response.Response{
		Code: response.StatusOK,
		Data: map[string]string{
			"token": res.Token,
		},
	})
}

func GetRandomFile(ctx *gin.Context) {
	token := ctx.GetHeader(common.FileTokenHeader)
	res, err := cli.Service.GetRandomFile(ctx, &pb.TokenReq{
		Token: token,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: response.StatusOK,
		Data: res.Files,
	})
}

func GetHotLikeRank(ctx *gin.Context) {
	token := ctx.GetHeader(common.FileTokenHeader)
	res, err := cli.Service.GetHotLikeRank(ctx, &pb.TokenReq{
		Token: token,
	})
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

func GetUserImages(ctx *gin.Context) {
	userName := ctx.Param("userName")
	if userName == "" {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: "check params",
		})
		return
	}
	var req pb.GetUserImagesReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: "check params",
		})
		return
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	token := ctx.GetHeader(common.FileTokenHeader)
	req.Token = token
	req.Userid = userName
	res, err := cli.Service.GetUserImages(ctx, &req)
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
			"total": res.Total,
			"files": res.Files,
		},
	})
}
