package main

import (
	"log"
	"time"

	"github.com/miRemid/kira/services/auth/client"
	"github.com/micro/go-micro/v2"
)

func main() {
	server := micro.NewService(
		micro.Name("kira.micro.client.auth"),
	)
	server.Init()

	cli := client.NewAuthClient(server)
	resp, err := cli.Auth("asdf", "asdfkdsjaf")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("TokenString = %s\n", resp.Token)

	resp2, err := cli.Valid(resp.Token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Request UserID: %s\n Valid UserID: %s\n", "asdf", resp2.UserId)

	log.Println("Waiting 5 seconds for freshing token test")
	time.Sleep(time.Second * 3)

	resp3, err := cli.Refresh(resp.Token)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("RefreshToken = %s\n", resp3.Token)
}
