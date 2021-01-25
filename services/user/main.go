package main

import (
	"log"

	hystrixGo "github.com/afex/hystrix-go/hystrix"
	"github.com/casbin/casbin/v2"
	"github.com/miRemid/kira/common"

	"github.com/miRemid/kira/services/user/handler"
	"github.com/miRemid/kira/services/user/pb"
	"github.com/miRemid/kira/services/user/repository"
	"github.com/miRemid/kira/services/user/route"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix/v2"
	"github.com/pkg/errors"
)

func startAPIService() {
	e, err := casbin.NewEnforcer("./casbin/model.conf", "./casbin/permission.csv")
	if err != nil {
		log.Fatal(err)
	}
	r := route.Route(e)
	service := web.NewService(
		web.Name("kira.micro.api.user"),
		web.Address(common.Getenv("API_ADDRESS", ":5002")),
		web.Handler(r),
		web.Registry(etcd.NewRegistry(
			registry.Addrs(common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")),
		)),
	)
	route.Init(client.DefaultClient)
	service.Init()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func startMicroService() {
	service := micro.NewService(
		micro.Name("kira.micro.service.user"),
		micro.Version("latest"),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")),
		)),
		micro.WrapClient(hystrix.NewClientWrapper()),
	)
	service.Init()
	hystrixGo.DefaultMaxConcurrent = 5
	hystrixGo.DefaultTimeout = 300

	db, err := common.DBConnect()
	if err != nil {
		log.Fatal(errors.WithMessage(err, "connect to database"))
	}

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

func main() {
	log.SetFlags(log.Llongfile)
	go startAPIService()
	startMicroService()
}
