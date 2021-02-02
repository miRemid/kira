package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/proto/pb"
)

type UserPassword struct {
	Username string `form:"user_name" binding:"required,usernameValidate"`
	Password string `form:"password" binding:"required,passwordValidate"`
}

func Signin(ctx *gin.Context) {
	var st UserPassword
	if err := ctx.ShouldBind(&st); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: err.Error(),
		})
		return
	}

	res, err := cli.Signin(st.Username, st.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    response.StatusOK,
		Message: res.Msg,
		Data:    res.Token,
	})

}

func Signup(ctx *gin.Context) {
	var st UserPassword
	if err := ctx.ShouldBind(&st); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: err.Error(),
		})
		return
	}
	res, err := cli.Signup(st.Username, st.Password)
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

func GetInfo(ctx *gin.Context) {
	userid := ctx.GetHeader("userid")

	res, err := cli.UserInfo(userid)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code:    response.StatusOK,
		Message: res.Msg,
		Data:    res.User,
	})
}

type DeleteReq struct {
	UserID string `json:"user_id" form:"user_id" binding:"required"`
}

func DeleteUser(ctx *gin.Context) {
	var req DeleteReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: err.Error(),
		})
		return
	}
	_, err := cli.DeleteUser(req.UserID)
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

func GetUserList(ctx *gin.Context) {
	var req pb.UserListRequest
	if err := ctx.BindQuery(&req); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: err.Error(),
		})
		return
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	res, err := cli.GetUserList(req.Limit, req.Offset)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code: response.StatusOK,
		Data: res,
	})
}

func UpdateUser(ctx *gin.Context) {
	var req pb.UpdateUserRoleRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusBadParams,
			Error: err.Error(),
		})
		return
	}
	res, err := cli.UpdateUser(req.UserID, req.Role)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:  response.StatusInternalError,
			Error: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code:    response.StatusOK,
		Message: res.Message,
	})
}
