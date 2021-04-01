package main

import (
	"log"

	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/wrapper/hystrix"
	"github.com/miRemid/kira/common/wrapper/tracer"
	"github.com/miRemid/kira/services/file-api/route"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"

	hystrixGo "github.com/afex/hystrix-go/hystrix"
)

func main() {
	log.SetFlags(log.Llongfile)
	etcdAddr := common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs(etcdAddr),
	)

	jaegerTracer, closer, err := tracer.NewJaegerTracer("kira.micro.client.file", common.Getenv("JAEGER_ADDRESS", "127.0.0.1:6831"))
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	r := route.Route()
	service := web.NewService(
		web.Name("kira.micro.api.file"),
		web.Address(common.Getenv("API_ADDRESS", ":5001")),
		web.Handler(r),
		web.Registry(etcdRegistry),
	)

	cli := micro.NewService(
		micro.Name("kira.micro.client.file"),
		micro.Registry(etcdRegistry),
		micro.WrapClient(
			hystrix.NewClientWrapper(),
			opentracing.NewClientWrapper(jaegerTracer),
		),
	)
	hystrixGo.DefaultSleepWindow = 5000
	hystrixGo.DefaultMaxConcurrent = 50

	route.Init(cli.Client())
	service.Init()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}