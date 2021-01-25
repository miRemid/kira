package client

import (
	"context"

	"github.com/miRemid/kira/services/auth/pb"
	"github.com/micro/go-micro/v2/client"
)

type AuthClient struct {
	service pb.AuthService
}

func NewAuthClient(client client.Client) *AuthClient {
	var cli AuthClient
	srv := pb.NewAuthService("kira.micro.service.auth", client)
	cli.service = srv
	return &cli
}

func (client *AuthClient) Auth(userid, userRole string) (*pb.AuthResponse, error) {
	return client.service.Auth(context.TODO(), &pb.AuthRequest{
		UserID:   userid,
		UserRole: userRole,
	})
}

func (client *AuthClient) Valid(tokenString string) (*pb.ValidResponse, error) {
	return client.service.Valid(context.TODO(), &pb.TokenRequest{
		Token: tokenString,
	})
}

func (client *AuthClient) Refresh(tokenString string) (*pb.AuthResponse, error) {
	return client.service.Refresh(context.TODO(), &pb.TokenRequest{
		Token: tokenString,
	})
}

func (client *AuthClient) FileToken(tokenString string) (*pb.FileTokenResponse, error) {
	return client.service.FileToken(context.TODO(), &pb.FileTokenRequest{
		Token: tokenString,
	})
}
