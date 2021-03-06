package handler

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/user/repository"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserHandler struct {
	Repo repository.UserRepository
}

func (handler UserHandler) ForgetPassword(ctx context.Context, in *pb.ForgetPasswordRequest, res *pb.ForgetPasswordResponse) error {
	err := handler.Repo.ForgetPassword(ctx, in.Username, in.Email)
	if err != nil {
		res.Error = err.Error()
		return err
	}
	res.Success = true
	return nil
}

func (handler UserHandler) ModifyPassword(ctx context.Context, in *pb.ModifyPasswordRequest, res *pb.ModifyPasswordResponse) error {
	code, err := handler.Repo.ModifyPassword(ctx, in.Random, in.Email, in.Password)
	if err != nil {
		res.Error = err.Error()
		return err
	}
	res.Success = true
	res.Code = code
	return nil
}

func (handler UserHandler) BindMail(ctx context.Context, in *pb.BindMailRequest, res *pb.BindMailResponse) error {
	err := handler.Repo.BindMail(ctx, in.Mail, in.Userid)
	if err != nil {
		res.Error = err.Error()
		return err
	}
	res.Success = true
	return nil
}

func (handler UserHandler) VertifyBindMail(ctx context.Context, in *pb.VertifyBindMailRequest, res *pb.VertifyBindMailResponse) error {
	code, err := handler.Repo.BindMailFinal(ctx, in.Random, in.Userid)
	if err != nil {
		res.Error = err.Error()
		return err
	}
	res.Success = true
	res.Code = int64(code)
	return nil
}

func (handler UserHandler) GetLoginUserInfo(ctx context.Context, in *pb.LoginUserInfoReq, res *pb.LoginUserInfoRes) error {
	user, token, err := handler.Repo.LoginUserInfo(ctx, in.Userid)
	res.User = new(pb.User)
	res.User.UserID = user.UserID
	res.User.UserName = user.UserName
	res.User.UserRole = user.Role
	res.User.UserStatus = user.Status
	res.Token = token
	return err
}

func (handler UserHandler) Signin(ctx context.Context, in *pb.SigninReq, res *pb.SigninRes) error {
	token, role, err := handler.Repo.Signin(ctx, in.Username, in.Password)
	if err != nil {
		res.Succ = false
		res.Msg = err.Error()
		return errors.WithMessage(err, "sign in")
	}
	res.Succ = true
	res.Msg = "sign in success"
	res.Token = token
	if role == "admin" {
		res.Admin = true
	} else {
		res.Admin = false
	}
	return nil
}

func (handler UserHandler) Signup(ctx context.Context, in *pb.SignupReq, res *pb.SignupRes) error {
	err := handler.Repo.Signup(ctx, in.Username, in.Password)
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
	user, err := handler.Repo.UserInfo(ctx, in.UserName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrRecordNotFound
		}
		return err
	}
	res.Succ = true
	res.Msg = "get user infomation success"
	res.User = new(pb.User)
	res.User.UserID = user.UserID
	res.User.UserName = user.UserName
	res.User.UserRole = user.Role
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
		u.Status = user.Status
		res.Users = append(res.Users, u)
	}
	return nil
}

func (handler UserHandler) AdminDeleteUser(ctx context.Context, in *pb.DeleteUserRequest, res *pb.AdminCommonResponse) error {
	return handler.Repo.DeleteUser(ctx, in.UserID)
}

func (handler UserHandler) AdminUpdateUser(ctx context.Context, in *pb.UpdateUserRoleRequest, res *pb.AdminCommonResponse) error {
	if err := handler.Repo.ChangeUserStatus(ctx, in.UserID, in.Status); err != nil {
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

func (handler UserHandler) ChangePassword(ctx context.Context, in *pb.UpdatePasswordReq, res *pb.UpdatePasswordRes) error {
	if err := handler.Repo.ChangePassword(ctx, in.UserID, in.OldPsw, in.NewPsw); err != nil {
		return err
	}
	res.Succ = true
	res.Msg = "change successful"
	return nil
}

func (handler UserHandler) GetUserImages(ctx context.Context, in *pb.GetUserImagesReqByNameReq, res *pb.GetUserImagesRes) error {
	resp, err := handler.Repo.GetUserImages(ctx, in.UserName, in.Offset, in.Limit, in.Desc)
	if err != nil {
		return errors.WithMessage(err, "Get User Images")
	}
	res.Files = resp.Files
	res.Total = resp.Total
	return nil
}
