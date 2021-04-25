package main

import (
	"log"

	hystrixGo "github.com/afex/hystrix-go/hystrix"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/wrapper/hystrix"
	"github.com/miRemid/kira/common/wrapper/tracer"
	"github.com/miRemid/kira/services/site/client"
	"github.com/miRemid/kira/services/site/router"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
)

func main() {
	log.SetFlags(log.Llongfile)

	etcdAddr := common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs(etcdAddr),
	)
	jaegerTracer, closer, err := tracer.NewJaegerTracer("kira.micro.client.site", common.Getenv("JAEGER_ADDRESS", "127.0.0.1:6831"))
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	cli := micro.NewService(
		micro.Name("kira.micro.client.site"),
		micro.Registry(etcdRegistry),
		micro.WrapClient(
			hystrix.NewClientWrapper(),
			opentracing.NewClientWrapper(jaegerTracer),
		),
	)
	cli.Init()

	hystrixGo.DefaultMaxConcurrent = 50
	hystrixGo.DefaultTimeout = 10000

	client.Init(cli.Client())
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
