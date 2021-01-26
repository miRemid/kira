package main

import (
	"log"
	"net/http"

	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/services/site/router"
	"github.com/pkg/errors"
)

func main() {
	log.SetFlags(log.Llongfile)
	r := router.New()

	if err := http.ListenAndServe(common.Getenv("API_ADDRESS", ":5000"), r); err != nil {
		log.Fatal(errors.WithMessage(err, "http service"))
	}
}
