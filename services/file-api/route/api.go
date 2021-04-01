package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/proto/pb"
)

type Search struct {
	Offset int64 `json:"offset" form:"offset"`
	Limit  int64 `json:"limit" form:"limit"`
}

type SearchRes struct {
	Total int64      `json:"total"`
	Files []*pb.File `json:"files"`
}

func GetHistory(ctx *gin.Context) {
	token, _ := ctx.Get("owner")
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

type DeleteReq struct {
	FileID string `json:"file_id" form:"file_id" binding:"required"`
}

func DeleteFile(ctx *gin.Context) {
	token, _ := ctx.Get("owner")
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
	if err != nil || !res.Succ {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	} else {
		ctx.JSON(http.StatusOK, response.Response{
			Code: response.StatusOK,
			Data: res.File,
		})
	}
}