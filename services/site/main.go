package main

import (
	"log"

	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/services/site/router"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
)

func main() {
	log.SetFlags(log.Llongfile)

	etcdAddr := common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs(etcdAddr),
	)

	r := router.New()

	service := web.NewService(
		web.Name("go.micro.api.site"),
		web.Address(common.Getenv("API_ADDRESS", ":5000")),
		web.Handler(r),
		web.Registry(etcdRegistry),
	)

	service.Init()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
