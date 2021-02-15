package main

import (
	"log"

	"github.com/dgrijalva/jwt-go/test"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/tracer"
	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/auth/handler"
	"github.com/miRemid/kira/services/auth/repository"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/pkg/errors"
)

func main() {
	log.SetFlags(log.Llongfile)

	jaegerTracer, closer, err := tracer.NewJaegerTracer("kira.micro.service.auth", common.Getenv("JAEGER_ADDRESS", "127.0.0.1:6831"))
	if err != nil {
		log.Fatal(errors.WithMessage(err, "tracer"))
	}
	defer closer.Close()

	db, err := common.DBConnect()
	if err != nil {
		log.Fatal(errors.WithMessage(err, "database"))
	}

	service := micro.NewService(
		micro.Name("kira.micro.service.auth"),
		micro.Version("latest"),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")),
		)),
		micro.WrapHandler(opentracing.NewHandlerWrapper(jaegerTracer)),
	)
	service.Init()

	prikey := test.LoadRSAPrivateKeyFromDisk(common.Getenv("PRIKEY_PATH", "./pem/prikey.pem"))
	pubkey := test.LoadRSAPublicKeyFromDisk(common.Getenv("PUBKEY_PATH", "./pem/pubkey.pem"))

	authHandler := &handler.AuthHandler{
		Repo: repository.NewAuthRepositoryImpl(pubkey, prikey, db),
	}
	if err := pb.RegisterAuthServiceHandler(service.Server(), authHandler); err != nil {
		log.Fatal(errors.WithMessage(err, "register server"))
	}

	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run server"))
	}
}
