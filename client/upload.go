package client

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"time"

	"github.com/miRemid/kira/proto/pb"
	"github.com/micro/go-micro/v2/client"
)

type UploadClient struct {
	Service pb.UploadService
}

func NewUploadClient(client client.Client) *UploadClient {
	var cli UploadClient
	srv := pb.NewUploadService("kira.micro.service.upload", client)
	cli.Service = srv
	return &cli
}

func (cli UploadClient) UploadFile(token, fileName, fileExt string, width, height string, anony bool, file multipart.File) (*pb.UploadFileRes, error) {
	var buf bytes.Buffer
	size, _ := io.Copy(&buf, file)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()
	return cli.Service.UploadFile(ctx, &pb.UploadFileReq{
		FileName: fileName,
		FileExt:  fileExt,
		FileBody: buf.Bytes(),
		FileSize: size,
		Token:    token,
		Width:    width,
		Height:   height,
		Anony:    anony,
	})
}
func (client *UploadClient) Ping() (*pb.Pong, error) {
	return client.Service.Ping(context.TODO(), &pb.Ping{})
}
