package main

import (
	"log"

	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/wrapper/tracer"
	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/user/handler"
	"github.com/miRemid/kira/services/user/publisher"
	"github.com/miRemid/kira/services/user/repository"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/broker/nats"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/pkg/errors"
)

func main() {
	log.SetFlags(log.Llongfile)

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
		micro.Broker(nats.NewBroker(
			broker.Addrs(common.Getenv("NATS_ADDRESS", "nats://127.0.0.1:4222")),
		)),
		micro.WrapHandler(opentracing.NewHandlerWrapper(jaegerTracer)),
	)
	service.Init()

	db, err := common.DBConnect()
	if err != nil {
		log.Fatal(errors.WithMessage(err, "connect to database"))
	}

	publisher.Init(service.Client())

	repo, err := repository.NewUserRepository(service.Client(), db)
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
