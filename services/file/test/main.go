package main

import (
	"log"

	"github.com/miRemid/kira/services/file/client"
	"github.com/micro/go-micro/v2"
)

const (
	userid = "zxykm"
	token  = "1mmJzTeI51bGYVS69f001rrD1tG"
)

var cli client.FileClient
var srv micro.Service

func TestCreateBucketAndToken() {
	res, err := cli.GenerateToken(userid)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Token=%v, UserID=%v\n", res.Token, userid)
}

func main() {
	srv = micro.NewService(
		micro.Name("kira.micro.client.file"),
	)
	srv.Init()
	cli = client.NewFileClient(srv.Client())

	res, err := cli.RefreshToken(userid)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("New Token = %s\n", res.Token)

}
