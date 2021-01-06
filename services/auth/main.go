package main

import (
	"log"

	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/services/auth/handler"
	"github.com/miRemid/kira/services/auth/pb"
	"github.com/miRemid/kira/services/auth/repository"
	"github.com/micro/go-micro/v2"
	"github.com/pkg/errors"
)

func main() {
	log.SetFlags(log.Llongfile)

	service := micro.NewService(
		micro.Name("kira.micro.service.auth"),
		micro.Version("latest"),
	)
	service.Init()

	authHandler := &handler.AuthHandler{
		Repo: repository.NewAuthRepositoryImpl(common.Getenv("screct", "kira")),
	}
	if err := pb.RegisterAuthServiceHandler(service.Server(), authHandler); err != nil {
		log.Fatal(errors.WithMessage(err, "register server"))
	}

	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run server"))
	}
}
