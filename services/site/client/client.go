package client

import (
	"github.com/miRemid/kira/client"
	"github.com/miRemid/kira/common"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
)

var (
	fileCli *client.FileClient
)

func init() {
	service := micro.NewService(
		micro.Name("kira.service.service.site"),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")),
		)),
	)
	service.Init()
	fileCli = client.NewFileClient(service.Client())
}

func File() *client.FileClient {
	return fileCli
}
