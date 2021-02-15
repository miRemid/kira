package main

import (
	"log"

	hystrixGo "github.com/afex/hystrix-go/hystrix"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/tracer"
	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/upload/handler"
	"github.com/miRemid/kira/services/upload/repository"
	"github.com/miRemid/kira/services/upload/router"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix/v2"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/pkg/errors"
)

func startMicroService() {
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
		micro.WrapClient(hystrix.NewClientWrapper()),
		micro.WrapHandler(opentracing.NewHandlerWrapper(jaegerTracer)),
	)
	service.Init()

	hystrixGo.DefaultMaxConcurrent = 50
	hystrixGo.DefaultTimeout = 3000

	db, err := common.DBConnect()
	if err != nil {

	}

	mini, err := common.MinioConnect()
	if err != nil {

	}

	repo := repository.NewRepository(db, mini)

	if err := pb.RegisterUploadServiceHandler(service.Server(), handler.Handler{
		Repo: repo,
	}); err != nil {

	}

	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run service"))
	}

}

func startAPIService() {
	cli := client.DefaultClient
	cli.Init(grpc.MaxSendMsgSize(5 * 1024 * 1024))
	r := router.NewRouter(client.DefaultClient)
	service := web.NewService(
		web.Name("kira.micro.api.upload"),
		web.Address(common.Getenv("API_ADDRESS", ":5003")),
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

func main() {
	log.SetFlags(log.Llongfile)
	go startAPIService()
	startMicroService()
}
