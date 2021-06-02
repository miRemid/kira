package main

import (
	"log"

	hystrixGo "github.com/afex/hystrix-go/hystrix"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/casbin"
	"github.com/miRemid/kira/common/wrapper/hystrix"
	"github.com/pkg/errors"

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

	db, err := common.DBConnect()
	if err != nil {
		log.Fatal(errors.WithMessage(err, "connect to database"))
	}
	e := casbin.New(db, "./casbin/model.conf")
	e.LoadPolicy()
	e.AddPolicy("normal", "/user/me", "GET")
	e.AddPolicy("normal", "/user/changePassword", "POST")
	e.AddPolicy("normal", "/user/deleteAccount", "DELETE")

	e.AddPolicy("admin", "/user/changePassword", "POST")
	e.AddPolicy("admin", "/user/admin/deleteUserFile", "DELETE")
	e.AddPolicy("admin", "/user/admin/getUserList", "GET")
	e.AddPolicy("admin", "/user/admin/updateUserStatus", "POST")
	e.AddPolicy("admin", "/user/admin/getAnonyList", "GET")
	e.AddPolicy("admin", "/user/admin/deleteAnony", "DELETE")
	e.SavePolicy()

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
	r := route.Route(e)

	hystrixGo.DefaultMaxConcurrent = 50
	hystrixGo.DefaultTimeout = 10000

	service := web.NewService(
		web.Name("go.micro.api.user"),
		web.Address(common.Getenv("API_ADDRESS", ":5002")),
		web.Handler(r),
		web.Registry(etcdRegistry),
	)
	service.Init()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
