package main

import (
	"log"

	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/upload/handler"
	"github.com/miRemid/kira/services/upload/repository"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/pkg/errors"
)

func main() {
	log.SetFlags(log.Llongfile)

	// jaegerTracer, closer, err := tracer.NewJaegerTracer("kira.micro.service.upload", common.Getenv("JAEGER_ADDRESS", "127.0.0.1:6831"))
	// if err != nil {
	// 	log.Fatal(errors.WithMessage(err, "tracer"))
	// }
	// defer closer.Close()

	service := micro.NewService(
		micro.Name("kira.micro.service.upload"),
		micro.Version("latest"),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")),
		)),
		// micro.WrapHandler(opentracing.NewHandlerWrapper(jaegerTracer)),
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

	repo := repository.NewRepository(db, mini)

	if err := pb.RegisterUploadServiceHandler(service.Server(), handler.Handler{
		Repo: repo,
	}); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run service"))
	}
}
