package client

import (
	"context"
	"time"

	"github.com/miRemid/kira/proto/pb"
	"github.com/micro/go-micro/v2/client"
)

type UserClient struct {
	service pb.UserService
}

func NewUserClient(client client.Client) *UserClient {
	var cli UserClient
	srv := pb.NewUserService("kira.micro.service.user", client)
	cli.service = srv
	return &cli
}

func (cli UserClient) Signup(username, password string) (*pb.SignupRes, error) {
	return cli.service.Signup(context.TODO(), &pb.SignupReq{
		Username: username,
		Password: password,
	})
}

func (cli UserClient) Signin(username, password string) (*pb.SigninRes, error) {
	return cli.service.Signin(context.TODO(), &pb.SigninReq{
		Username: username,
		Password: password,
	})
}

func (cli UserClient) UserInfo(userid string) (*pb.UserInfoRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return cli.service.UserInfo(ctx, &pb.UserInfoReq{
		UserID: userid,
	})
}

func (cli UserClient) DeleteUser(userid string) (*pb.AdminCommonResponse, error) {
	return cli.service.AdminDeleteUser(context.TODO(), &pb.DeleteUserRequest{
		UserID: userid,
	})
}

func (cli UserClient) UpdateUser(userid string, status int64) (*pb.AdminCommonResponse, error) {
	return cli.service.AdminUpdateUser(context.TODO(), &pb.UpdateUserRoleRequest{
		UserID: userid,
		Status: status,
	})
}

func (cli UserClient) GetUserList(limit, offset int64) (*pb.UserListResponse, error) {
	return cli.service.AdminUserList(context.TODO(), &pb.UserListRequest{
		Limit:  limit,
		Offset: offset,
	})
}

func (client *UserClient) Ping() (*pb.Pong, error) {
	return client.service.Ping(context.TODO(), &pb.Ping{})
}

func (client *UserClient) ChangePassword(req *pb.UpdatePasswordReq) (*pb.UpdatePasswordRes, error) {
	return client.service.ChangePassword(context.TODO(), req)
}
