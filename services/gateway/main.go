package main

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/miRemid/kira/services/gateway/plugins/auth"

	"github.com/micro/micro/v2/client/api"
	"github.com/micro/micro/v2/cmd"
	"github.com/pkg/errors"
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	e, err := casbin.NewEnforcer("./casbin/model.conf", "./casbin/permission.csv")
	if err != nil {
		log.Fatal(err)
	}

	err = api.Register(auth.NewPlugin(e))
	if err != nil {
		log.Fatal(errors.WithMessage(err, "auth register"))
	}

	cmd.Init()
}
