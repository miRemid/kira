package main

import (
	"log"

	hystrixGo "github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/broker/nats"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/pkg/errors"

	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/wrapper/hystrix"
	"github.com/miRemid/kira/common/wrapper/tracer"
	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/file/handler"
	"github.com/miRemid/kira/services/file/repository"
	"github.com/miRemid/kira/services/file/route"
)

func startAPIService() {
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

func startMicroService() {
	jaegerTracer, closer, err := tracer.NewJaegerTracer("kira.micro.service.file", common.Getenv("JAEGER_ADDRESS", "127.0.0.1:6831"))
	if err != nil {
		log.Fatal(errors.WithMessage(err, "tracer"))
	}
	defer closer.Close()

	service := micro.NewService(
		micro.Name("kira.micro.service.file"),
		micro.Version("latest"),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")),
		)),
		micro.WrapHandler(opentracing.NewHandlerWrapper(jaegerTracer)),
		micro.Broker(nats.NewBroker(
			broker.Addrs(common.Getenv("NATS_ADDRESS", "nats://127.0.0.1:4222")),
		)),
	)
	service.Init()

	db, err := common.DBConnect()
	if err != nil {
		log.Fatal(errors.WithMessage(err, "connect to database"))
	}

	mini, err := common.MinioConnect()
	if err != nil {
		log.Fatal(errors.WithMessage(err, "connect to minio"))
	}

	repo := repository.NewFileRepository(db, mini)
	fileHandler := handler.FileServiceHandler{
		Repo: repo,
	}
	if err := pb.RegisterFileServiceHandler(service.Server(), fileHandler); err != nil {
		log.Fatal(errors.WithMessage(err, "register service"))
	}

	// 订阅消费者
	if err := micro.RegisterSubscriber("kira.micro.service.user.delete", service.Server(), fileHandler.DeleteUser); err != nil {
		log.Fatal(errors.WithMessage(err, "register subscriber"))
	}

	if err := service.Run(); err != nil {
		repo.Done()
		log.Fatal(errors.WithMessage(err, "run service"))
	}
}

func main() {
	log.SetFlags(log.Llongfile)
	go startAPIService()
	startMicroService()
}
