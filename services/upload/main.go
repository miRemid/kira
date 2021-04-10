package main

import (
	"log"

	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/wrapper/tracer"
	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/upload/handler"
	"github.com/miRemid/kira/services/upload/repository"
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

	jaegerTracer, closer, err := tracer.NewJaegerTracer("kira.micro.service.upload", common.Getenv("JAEGER_ADDRESS", "127.0.0.1:6831"))
	if err != nil {
		log.Fatal(errors.WithMessage(err, "tracer"))
	}
	defer closer.Close()

	service := micro.NewService(
		micro.Name("kira.micro.service.upload"),
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
		log.Fatal(err)
	}

	mini, err := common.MinioConnect()
	if err != nil {
		log.Fatal(err)
	}

	publisher := micro.NewPublisher(common.AnonyEvent, service.Client())

	repo := repository.NewRepository(db, mini, publisher)

	if err := pb.RegisterUploadServiceHandler(service.Server(), handler.Handler{
		Repo: repo,
	}); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run service"))
	}
}
