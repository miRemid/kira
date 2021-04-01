package main

import (
	"log"

	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/wrapper/tracer"
	"github.com/miRemid/kira/services/gateway/plugins/auth"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/micro/micro/v2/client/api"
	"github.com/micro/micro/v2/cmd"
)

func main() {
	log.SetFlags(log.Llongfile)

	etcdAddr := common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs(etcdAddr),
	)

	jaegerTracer, closer, err := tracer.NewJaegerTracer("kira.micro.client.gateway", common.Getenv("JAEGER_ADDRESS", "127.0.0.1:6831"))
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	cli := micro.NewService(
		micro.Name("kira.micro.client.gateway"),
		micro.Registry(etcdRegistry),
		micro.WrapClient(
			opentracing.NewClientWrapper(jaegerTracer),
		),
	)

	err = api.Register(auth.NewPlugin(cli.Client()))
	if err != nil {
		log.Fatal("auth register")
	}

	cmd.Init()
}
