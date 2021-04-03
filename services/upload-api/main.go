package main

import (
	"log"

	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/wrapper/tracer"
	"github.com/miRemid/kira/services/upload-api/router"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
)

func main() {
	log.SetFlags(log.Llongfile)

	etcdRegistry := etcd.NewRegistry(
		registry.Addrs(common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")),
	)

	jaegerTracer, closer, err := tracer.NewJaegerTracer("kira.micro.client.upload", common.Getenv("JAEGER_ADDRESS", "127.0.0.1:6831"))
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	cli := micro.NewService(
		micro.Name("kira.micro.client.upload"),
		micro.Registry(etcdRegistry),
		micro.WrapClient(
			// hystrix.NewClientWrapper(),
			opentracing.NewClientWrapper(jaegerTracer),
		),
	)

	// hystrixGo.DefaultMaxConcurrent = 50
	// hystrixGo.DefaultTimeout = 10000

	cli.Client().Init(grpc.MaxSendMsgSize(5 * 1024 * 1024))
	r := router.NewRouter(cli.Client())
	service := web.NewService(
		web.Name("go.micro.api.upload"),
		web.Address(common.Getenv("API_ADDRESS", ":5003")),
		web.Handler(r),
		web.Registry(etcdRegistry),
	)
	service.Init()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
