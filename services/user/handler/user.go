package handler

import (
	"context"

	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/user/repository"
	"github.com/pkg/errors"
)

type UserHandler struct {
	Repo repository.UserRepository
}

func (handler UserHandler) Signin(ctx context.Context, in *pb.SigninReq, res *pb.SigninRes) error {
	token, err := handler.Repo.Signin(in.Username, in.Password)
	if err != nil {
		res.Succ = false
		res.Msg = err.Error()
		return errors.WithMessage(err, "sign in")
	}
	res.Succ = true
	res.Msg = "sign in success"
	res.Token = token
	return nil
}

func (handler UserHandler) Signup(ctx context.Context, in *pb.SignupReq, res *pb.SignupRes) error {
	err := handler.Repo.Signup(in.Username, in.Password)
	if err != nil {
		res.Succ = false
		res.Msg = err.Error()
		return errors.WithMessage(err, "sign up")
	}
	res.Succ = true
	res.Msg = "sign up success"
	return nil
}
func (handler UserHandler) UserInfo(ctx context.Context, in *pb.UserInfoReq, res *pb.UserInfoRes) error {
	user, err := handler.Repo.UserInfo(in.UserID)
	if err != nil {
		res.Succ = false
		res.Msg = err.Error()
		return errors.WithMessage(err, "get user infomation")
	}
	res.Succ = true
	res.Msg = "get user infomation success"
	res.User = new(pb.User)
	res.User.UserID = user.UserID
	res.User.UserName = user.UserName
	res.User.UserRole = user.Role
	res.User.UserToken = user.Token
	return nil
}

func (handler UserHandler) AdminUserList(ctx context.Context, in *pb.UserListRequest, res *pb.UserListResponse) error {
	return nil
}

func (handler UserHandler) AdminDeleteUser(ctx context.Context, in *pb.DeleteUserRequest, res *pb.AdminCommonResponse) error {
	return handler.Repo.DeleteUser(ctx, in.UserID)
}
func (handler UserHandler) AdminUpdateUser(ctx context.Context, in *pb.UpdateUserRoleRequest, res *pb.AdminCommonResponse) error {
	return nil
}
