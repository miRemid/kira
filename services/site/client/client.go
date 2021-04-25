package client

import (
	"github.com/miRemid/kira/client"
	mClient "github.com/micro/go-micro/v2/client"
)

var (
	FileCli   *client.FileClient
	UserCli   *client.UserClient
	AuthCli   *client.AuthClient
	UploadCli *client.UploadClient
)

func Init(cli mClient.Client) {
	FileCli = client.NewFileClient(cli)
	UserCli = client.NewUserClient(cli)
	AuthCli = client.NewAuthClient(cli)
	UploadCli = client.NewUploadClient(cli)
}

func File() *client.FileClient {
	return FileCli
}
