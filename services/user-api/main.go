package main

import (
	"log"

	hystrixGo "github.com/afex/hystrix-go/hystrix"
	"github.com/casbin/casbin/v2"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/wrapper/hystrix"

	"github.com/miRemid/kira/common/wrapper/tracer"
	"github.com/miRemid/kira/services/user-api/route"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
)

func main() {
	log.SetFlags(log.Llongfile)

	e, err := casbin.NewEnforcer("./casbin/model.conf", "./casbin/permission.csv")
	if err != nil {
		log.Fatal(err)
	}
	r := route.Route(e)

	etcdAddr := common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs(etcdAddr),
	)

	jaegerTracer, closer, err := tracer.NewJaegerTracer("kira.micro.client.user", common.Getenv("JAEGER_ADDRESS", "127.0.0.1:6831"))
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	cli := micro.NewService(
		micro.Name("kira.micro.client.user"),
		micro.Registry(etcdRegistry),
		micro.WrapClient(
			hystrix.NewClientWrapper(),
			opentracing.NewClientWrapper(jaegerTracer),
		),
	)

	route.Init(cli.Client())

	hystrixGo.DefaultMaxConcurrent = 50
	hystrixGo.DefaultTimeout = 10000

	service := web.NewService(
		web.Name("go.micro.api.user"),
		web.Address(common.Getenv("API_ADDRESS", ":5002")),
		web.Handler(r),
		web.Registry(etcd.NewRegistry(
			registry.Addrs(common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")),
		)),
	)
	service.Init()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
