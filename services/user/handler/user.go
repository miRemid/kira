package handler

import (
	"context"

	"github.com/golang/protobuf/ptypes"
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
	res.User.UserStatus = user.Status
	return nil
}

func (handler UserHandler) AdminUserList(ctx context.Context, in *pb.UserListRequest, res *pb.UserListResponse) error {
	users, total, err := handler.Repo.GetUserList(ctx, in.Limit, in.Offset)
	if err != nil {
		return errors.WithMessage(err, "get user list")
	}
	res.Total = total
	res.Users = make([]*pb.UserListResponse_User, 0)
	for _, user := range users {
		u := new(pb.UserListResponse_User)
		u.CreateTime, _ = ptypes.TimestampProto(user.CreatedAt)
		u.Role = user.Role
		u.UserID = user.UserID
		u.UserName = user.UserName
		res.Users = append(res.Users, u)
	}
	return nil
}

func (handler UserHandler) AdminDeleteUser(ctx context.Context, in *pb.DeleteUserRequest, res *pb.AdminCommonResponse) error {
	return handler.Repo.DeleteUser(ctx, in.UserID)
}
func (handler UserHandler) AdminUpdateUser(ctx context.Context, in *pb.UpdateUserRoleRequest, res *pb.AdminCommonResponse) error {
	if err := handler.Repo.ChangeUserRole(ctx, in.UserID, in.Role); err != nil {
		return err
	}
	res.Message = "update success"
	return nil
}

func (handler UserHandler) Ping(ctx context.Context, in *pb.Ping, res *pb.Pong) error {
	res.Code = 0
	res.Name = "user"
	res.Message = "ok"
	return nil
}
