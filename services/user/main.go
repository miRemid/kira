package main

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/tracer"

	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/user/handler"
	"github.com/miRemid/kira/services/user/repository"
	"github.com/miRemid/kira/services/user/route"

	hystrixGo "github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/broker/nats"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix/v2"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/pkg/errors"
)

func startAPIService() {
	e, err := casbin.NewEnforcer("./casbin/model.conf", "./casbin/permission.csv")
	if err != nil {
		log.Fatal(err)
	}
	r := route.Route(e)
	route.Init(client.DefaultClient)

	service := web.NewService(
		web.Name("kira.micro.api.user"),
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

func startMicroService() {
	jaegerTracer, closer, err := tracer.NewJaegerTracer("kira.micro.service.user", common.Getenv("JAEGER_ADDRESS", "127.0.0.1:6831"))
	if err != nil {
		log.Fatal(errors.WithMessage(err, "tracer"))
	}
	defer closer.Close()
	service := micro.NewService(
		micro.Name("kira.micro.service.user"),
		micro.Version("latest"),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")),
		)),
		micro.WrapClient(hystrix.NewClientWrapper()),
		micro.Broker(nats.NewBroker(
			broker.Addrs(common.Getenv("NATS_ADDRESS", "nats://127.0.0.1:4222")),
		)),
		micro.WrapHandler(opentracing.NewHandlerWrapper(jaegerTracer)),
	)
	service.Init()
	hystrixGo.DefaultMaxConcurrent = 50
	hystrixGo.DefaultTimeout = 5000

	db, err := common.DBConnect()
	if err != nil {
		log.Fatal(errors.WithMessage(err, "connect to database"))
	}

	pub := micro.NewPublisher("kira.micro.service.user.delete", service.Client())

	repo, err := repository.NewUserRepository(service.Client(), db, pub)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "new repo"))
	}

	userHandler := handler.UserHandler{
		Repo: repo,
	}

	if err := pb.RegisterUserServiceHandler(service.Server(), userHandler); err != nil {
		log.Fatal(errors.WithMessage(err, "register handler"))
	}

	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run service"))
	}
}

func main() {
	log.SetFlags(log.Llongfile)
	go startAPIService()
	startMicroService()
}
