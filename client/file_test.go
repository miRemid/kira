package client_test

import (
	"testing"

	"github.com/miRemid/kira/client"
	"github.com/miRemid/kira/common"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
)

func TestFile(t *testing.T) {
	service := micro.NewService(
		micro.Name("kira.micro.test.file"),
		micro.Version("latest"),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")),
		)),
	)
	service.Init()
	cli := client.NewFileClient(service.Client())
	res, err := cli.GetRandomFile()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res.Files)
}
