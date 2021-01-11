package main

import (
	"log"
	"time"

	"github.com/miRemid/kira/services/auth/client"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
)

func main() {
	server := micro.NewService(
		micro.Name("kira.micro.client.auth"),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs("127.0.0.1:2379"),
		)),
	)
	server.Init()

	cli := client.NewAuthClient(server.Client())
	resp, err := cli.Auth("asdf", "asdfkdsjaf")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("TokenString = %s\n", resp.Token)

	resp2, err := cli.Valid(resp.Token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Request UserID: %s\n Valid UserID: %s\n", "asdf", resp2.UserID)

	log.Println("Waiting 5 seconds for freshing token test")
	time.Sleep(time.Second * 3)

	resp3, err := cli.Refresh(resp.Token)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("RefreshToken = %s\n", resp3.Token)
}
