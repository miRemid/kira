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
	service pb.UploadService
}

func NewUploadClient(client client.Client) *UploadClient {
	var cli UploadClient
	srv := pb.NewUploadService("kira.micro.service.upload", client)
	cli.service = srv
	return &cli
}

func (cli UploadClient) UploadFile(owner, fileName, fileExt string, width, height string, anony bool, file multipart.File) (*pb.UploadFileRes, error) {
	var buf bytes.Buffer
	size, _ := io.Copy(&buf, file)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()
	return cli.service.UploadFile(ctx, &pb.UploadFileReq{
		FileName: fileName,
		FileExt:  fileExt,
		FileBody: buf.Bytes(),
		FileSize: size,
		Owner:    owner,
		Width:    width,
		Height:   height,
		Anony:    anony,
	})
}
func (client *UploadClient) Ping() (*pb.Pong, error) {
	return client.service.Ping(context.TODO(), &pb.Ping{})
}
