// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/user.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for UserService service

func NewUserServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for UserService service

type UserService interface {
	// API
	// Common
	Signin(ctx context.Context, in *SigninReq, opts ...client.CallOption) (*SigninRes, error)
	Signup(ctx context.Context, in *SignupReq, opts ...client.CallOption) (*SignupRes, error)
	GetUserImages(ctx context.Context, in *GetUserImagesReqByNameReq, opts ...client.CallOption) (*GetUserImagesRes, error)
	UserInfo(ctx context.Context, in *UserInfoReq, opts ...client.CallOption) (*UserInfoRes, error)
	ForgetPassword(ctx context.Context, in *ForgetPasswordRequest, opts ...client.CallOption) (*ForgetPasswordResponse, error)
	ModifyPassword(ctx context.Context, in *ModifyPasswordRequest, opts ...client.CallOption) (*ModifyPasswordResponse, error)
	// User
	ChangePassword(ctx context.Context, in *UpdatePasswordReq, opts ...client.CallOption) (*UpdatePasswordRes, error)
	GetLoginUserInfo(ctx context.Context, in *LoginUserInfoReq, opts ...client.CallOption) (*LoginUserInfoRes, error)
	BindMail(ctx context.Context, in *BindMailRequest, opts ...client.CallOption) (*BindMailResponse, error)
	VertifyBindMail(ctx context.Context, in *VertifyBindMailRequest, opts ...client.CallOption) (*VertifyBindMailResponse, error)
	// Admin
	AdminUserList(ctx context.Context, in *UserListRequest, opts ...client.CallOption) (*UserListResponse, error)
	AdminDeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...client.CallOption) (*AdminCommonResponse, error)
	AdminUpdateUser(ctx context.Context, in *UpdateUserRoleRequest, opts ...client.CallOption) (*AdminCommonResponse, error)
	// RPC
	Ping(ctx context.Context, in *Ping, opts ...client.CallOption) (*Pong, error)
}

type userService struct {
	c    client.Client
	name string
}

func NewUserService(name string, c client.Client) UserService {
	return &userService{
		c:    c,
		name: name,
	}
}

func (c *userService) Signin(ctx context.Context, in *SigninReq, opts ...client.CallOption) (*SigninRes, error) {
	req := c.c.NewRequest(c.name, "UserService.Signin", in)
	out := new(SigninRes)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) Signup(ctx context.Context, in *SignupReq, opts ...client.CallOption) (*SignupRes, error) {
	req := c.c.NewRequest(c.name, "UserService.Signup", in)
	out := new(SignupRes)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) GetUserImages(ctx context.Context, in *GetUserImagesReqByNameReq, opts ...client.CallOption) (*GetUserImagesRes, error) {
	req := c.c.NewRequest(c.name, "UserService.GetUserImages", in)
	out := new(GetUserImagesRes)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) UserInfo(ctx context.Context, in *UserInfoReq, opts ...client.CallOption) (*UserInfoRes, error) {
	req := c.c.NewRequest(c.name, "UserService.UserInfo", in)
	out := new(UserInfoRes)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) ForgetPassword(ctx context.Context, in *ForgetPasswordRequest, opts ...client.CallOption) (*ForgetPasswordResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.ForgetPassword", in)
	out := new(ForgetPasswordResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) ModifyPassword(ctx context.Context, in *ModifyPasswordRequest, opts ...client.CallOption) (*ModifyPasswordResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.ModifyPassword", in)
	out := new(ModifyPasswordResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) ChangePassword(ctx context.Context, in *UpdatePasswordReq, opts ...client.CallOption) (*UpdatePasswordRes, error) {
	req := c.c.NewRequest(c.name, "UserService.ChangePassword", in)
	out := new(UpdatePasswordRes)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) GetLoginUserInfo(ctx context.Context, in *LoginUserInfoReq, opts ...client.CallOption) (*LoginUserInfoRes, error) {
	req := c.c.NewRequest(c.name, "UserService.GetLoginUserInfo", in)
	out := new(LoginUserInfoRes)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) BindMail(ctx context.Context, in *BindMailRequest, opts ...client.CallOption) (*BindMailResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.BindMail", in)
	out := new(BindMailResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) VertifyBindMail(ctx context.Context, in *VertifyBindMailRequest, opts ...client.CallOption) (*VertifyBindMailResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.VertifyBindMail", in)
	out := new(VertifyBindMailResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) AdminUserList(ctx context.Context, in *UserListRequest, opts ...client.CallOption) (*UserListResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.AdminUserList", in)
	out := new(UserListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) AdminDeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...client.CallOption) (*AdminCommonResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.AdminDeleteUser", in)
	out := new(AdminCommonResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) AdminUpdateUser(ctx context.Context, in *UpdateUserRoleRequest, opts ...client.CallOption) (*AdminCommonResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.AdminUpdateUser", in)
	out := new(AdminCommonResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) Ping(ctx context.Context, in *Ping, opts ...client.CallOption) (*Pong, error) {
	req := c.c.NewRequest(c.name, "UserService.Ping", in)
	out := new(Pong)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserService service

type UserServiceHandler interface {
	// API
	// Common
	Signin(context.Context, *SigninReq, *SigninRes) error
	Signup(context.Context, *SignupReq, *SignupRes) error
	GetUserImages(context.Context, *GetUserImagesReqByNameReq, *GetUserImagesRes) error
	UserInfo(context.Context, *UserInfoReq, *UserInfoRes) error
	ForgetPassword(context.Context, *ForgetPasswordRequest, *ForgetPasswordResponse) error
	ModifyPassword(context.Context, *ModifyPasswordRequest, *ModifyPasswordResponse) error
	// User
	ChangePassword(context.Context, *UpdatePasswordReq, *UpdatePasswordRes) error
	GetLoginUserInfo(context.Context, *LoginUserInfoReq, *LoginUserInfoRes) error
	BindMail(context.Context, *BindMailRequest, *BindMailResponse) error
	VertifyBindMail(context.Context, *VertifyBindMailRequest, *VertifyBindMailResponse) error
	// Admin
	AdminUserList(context.Context, *UserListRequest, *UserListResponse) error
	AdminDeleteUser(context.Context, *DeleteUserRequest, *AdminCommonResponse) error
	AdminUpdateUser(context.Context, *UpdateUserRoleRequest, *AdminCommonResponse) error
	// RPC
	Ping(context.Context, *Ping, *Pong) error
}

func RegisterUserServiceHandler(s server.Server, hdlr UserServiceHandler, opts ...server.HandlerOption) error {
	type userService interface {
		Signin(ctx context.Context, in *SigninReq, out *SigninRes) error
		Signup(ctx context.Context, in *SignupReq, out *SignupRes) error
		GetUserImages(ctx context.Context, in *GetUserImagesReqByNameReq, out *GetUserImagesRes) error
		UserInfo(ctx context.Context, in *UserInfoReq, out *UserInfoRes) error
		ForgetPassword(ctx context.Context, in *ForgetPasswordRequest, out *ForgetPasswordResponse) error
		ModifyPassword(ctx context.Context, in *ModifyPasswordRequest, out *ModifyPasswordResponse) error
		ChangePassword(ctx context.Context, in *UpdatePasswordReq, out *UpdatePasswordRes) error
		GetLoginUserInfo(ctx context.Context, in *LoginUserInfoReq, out *LoginUserInfoRes) error
		BindMail(ctx context.Context, in *BindMailRequest, out *BindMailResponse) error
		VertifyBindMail(ctx context.Context, in *VertifyBindMailRequest, out *VertifyBindMailResponse) error
		AdminUserList(ctx context.Context, in *UserListRequest, out *UserListResponse) error
		AdminDeleteUser(ctx context.Context, in *DeleteUserRequest, out *AdminCommonResponse) error
		AdminUpdateUser(ctx context.Context, in *UpdateUserRoleRequest, out *AdminCommonResponse) error
		Ping(ctx context.Context, in *Ping, out *Pong) error
	}
	type UserService struct {
		userService
	}
	h := &userServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&UserService{h}, opts...))
}

type userServiceHandler struct {
	UserServiceHandler
}

func (h *userServiceHandler) Signin(ctx context.Context, in *SigninReq, out *SigninRes) error {
	return h.UserServiceHandler.Signin(ctx, in, out)
}

func (h *userServiceHandler) Signup(ctx context.Context, in *SignupReq, out *SignupRes) error {
	return h.UserServiceHandler.Signup(ctx, in, out)
}

func (h *userServiceHandler) GetUserImages(ctx context.Context, in *GetUserImagesReqByNameReq, out *GetUserImagesRes) error {
	return h.UserServiceHandler.GetUserImages(ctx, in, out)
}

func (h *userServiceHandler) UserInfo(ctx context.Context, in *UserInfoReq, out *UserInfoRes) error {
	return h.UserServiceHandler.UserInfo(ctx, in, out)
}

func (h *userServiceHandler) ForgetPassword(ctx context.Context, in *ForgetPasswordRequest, out *ForgetPasswordResponse) error {
	return h.UserServiceHandler.ForgetPassword(ctx, in, out)
}

func (h *userServiceHandler) ModifyPassword(ctx context.Context, in *ModifyPasswordRequest, out *ModifyPasswordResponse) error {
	return h.UserServiceHandler.ModifyPassword(ctx, in, out)
}

func (h *userServiceHandler) ChangePassword(ctx context.Context, in *UpdatePasswordReq, out *UpdatePasswordRes) error {
	return h.UserServiceHandler.ChangePassword(ctx, in, out)
}

func (h *userServiceHandler) GetLoginUserInfo(ctx context.Context, in *LoginUserInfoReq, out *LoginUserInfoRes) error {
	return h.UserServiceHandler.GetLoginUserInfo(ctx, in, out)
}

func (h *userServiceHandler) BindMail(ctx context.Context, in *BindMailRequest, out *BindMailResponse) error {
	return h.UserServiceHandler.BindMail(ctx, in, out)
}

func (h *userServiceHandler) VertifyBindMail(ctx context.Context, in *VertifyBindMailRequest, out *VertifyBindMailResponse) error {
	return h.UserServiceHandler.VertifyBindMail(ctx, in, out)
}

func (h *userServiceHandler) AdminUserList(ctx context.Context, in *UserListRequest, out *UserListResponse) error {
	return h.UserServiceHandler.AdminUserList(ctx, in, out)
}

func (h *userServiceHandler) AdminDeleteUser(ctx context.Context, in *DeleteUserRequest, out *AdminCommonResponse) error {
	return h.UserServiceHandler.AdminDeleteUser(ctx, in, out)
}

func (h *userServiceHandler) AdminUpdateUser(ctx context.Context, in *UpdateUserRoleRequest, out *AdminCommonResponse) error {
	return h.UserServiceHandler.AdminUpdateUser(ctx, in, out)
}

func (h *userServiceHandler) Ping(ctx context.Context, in *Ping, out *Pong) error {
	return h.UserServiceHandler.Ping(ctx, in, out)
}
