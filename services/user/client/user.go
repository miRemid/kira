package client

import (
	"context"

	"github.com/miRemid/kira/services/user/pb"
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
	return cli.service.UserInfo(context.TODO(), &pb.UserInfoReq{
		UserID: userid,
	})
}

func (cli UserClient) Refresh(userid string) (*pb.UserTokenRes, error) {
	return cli.service.RefreshToken(context.TODO(), &pb.UserTokenReq{
		UserID: userid,
	})
}
