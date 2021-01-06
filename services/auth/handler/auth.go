// handler 实现micro接口
package handler

import (
	"context"

	"github.com/miRemid/kira/services/auth/pb"
	"github.com/miRemid/kira/services/auth/repository"
	"github.com/miRemid/kira/services/auth/token"
	"github.com/pkg/errors"
)

type AuthHandler struct {
	Repo repository.AuthRepository
}

func (handler AuthHandler) Auth(ctx context.Context, in *pb.AuthRequest, res *pb.AuthResponse) error {
	tokenString, err := handler.Repo.Auth(ctx, in.Userid, in.Username)
	if err != nil {
		res.Succ = false
		res.Msg = err.Error()
		return errors.WithMessage(err, "generate token")
	}
	res.Succ = true
	res.Msg = "generate success"
	res.Token = tokenString
	return nil
}

func (handler AuthHandler) Valid(ctx context.Context, in *pb.TokenRequest, res *pb.ValidResponse) error {
	tk, err := handler.Repo.Valid(ctx, in.Token)
	if err != nil {
		res.Succ = false
		res.Msg = err.Error()
		return errors.WithMessage(err, "valid error")
	}
	res.Succ = true
	res.Msg = "valid token"
	claims := tk.Claims.(*token.AuthClaims)
	res.UserId = claims.UserID
	res.UserRole = claims.Role
	return nil
}

func (handler AuthHandler) Refresh(ctx context.Context, in *pb.TokenRequest, res *pb.AuthResponse) error {
	tokenString, err := handler.Repo.Refresh(ctx, in.Token)
	if err != nil {
		res.Succ = false
		res.Msg = err.Error()
		return errors.WithMessage(err, "generate token")
	}
	res.Succ = true
	res.Msg = "generate success"
	res.Token = tokenString
	return nil
}
