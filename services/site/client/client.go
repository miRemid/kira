package client

import (
	fileClient "github.com/miRemid/kira/services/file/client"
	"github.com/micro/go-micro/v2"
)

var (
	fileCli *fileClient.FileClient
)

func init() {
	service := micro.NewService(
		micro.Name("kira.service.client.site"),
	)
	service.Init()
	fileCli = fileClient.NewFileClient(service.Client())
}

func File() *fileClient.FileClient {
	return fileCli
}
