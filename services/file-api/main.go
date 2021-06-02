package main

import (
	"log"

	hystrixGo "github.com/afex/hystrix-go/hystrix"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/casbin"
	"github.com/pkg/errors"

	"github.com/miRemid/kira/common/wrapper/hystrix"
	"github.com/miRemid/kira/common/wrapper/tracer"
	"github.com/miRemid/kira/services/file-api/route"
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

	jaegerTracer, closer, err := tracer.NewJaegerTracer("kira.micro.client.file", common.Getenv("JAEGER_ADDRESS", "127.0.0.1:6831"))
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	db, err := common.DBConnect()
	if err != nil {
		log.Fatal(errors.WithMessage(err, "connect to database"))
	}
	e := casbin.New(db, "./casbin/model.conf")
	e.LoadPolicy()
	// e.AddPolicy("normal", "/file/getToken", "GET")
	// e.AddPolicy("normal", "/file/like", "POST")
	// e.AddPolicy("normal", "/file/getLikes", "GET")
	// e.SavePolicy()

	cli := micro.NewService(
		micro.Name("kira.micro.client.file"),
		micro.Registry(etcdRegistry),
		micro.WrapClient(
			hystrix.NewClientWrapper(),
			opentracing.NewClientWrapper(jaegerTracer),
		),
	)
	route.Init(cli.Client())
	r := route.Route(e)

	hystrixGo.DefaultMaxConcurrent = 50
	hystrixGo.DefaultTimeout = 10000

	service := web.NewService(
		web.Name("go.micro.api.file"),
		web.Address(common.Getenv("API_ADDRESS", ":5001")),
		web.Handler(r),
		web.Registry(etcdRegistry),
	)
	service.Init()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
