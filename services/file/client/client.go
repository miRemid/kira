package client

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"

	"github.com/miRemid/kira/services/file/pb"
	"github.com/micro/go-micro/v2/client"
)

type FileClient struct {
	service pb.FileService
}

func NewFileClient(service client.Client) FileClient {
	var cli FileClient
	srv := pb.NewFileService("kira.micro.service.file", service)
	cli.service = srv
	return cli
}

func (cli FileClient) GenerateToken(userid string) (*pb.TokenUserRes, error) {
	return cli.service.GenerateToken(context.TODO(), &pb.TokenUserReq{
		Userid: userid,
	})
}

func (cli FileClient) RefreshToken(userid string) (*pb.TokenUserRes, error) {
	return cli.service.RefreshToken(context.TODO(), &pb.TokenUserReq{
		Userid: userid,
	})
}

func (cli FileClient) GetHistory(token string, limit, offset int64) (*pb.GetHistoryRes, error) {
	return cli.service.GetHistory(context.TODO(), &pb.GetHistoryReq{
		Token:  token,
		Limit:  limit,
		Offset: offset,
	})
}

func (cli FileClient) UploadFile(token string, fileName, fileExt string, file multipart.File) (*pb.UploadFileRes, error) {
	var buf bytes.Buffer
	size, _ := io.Copy(&buf, file)
	return cli.service.UploadFile(context.TODO(), &pb.UploadFileReq{
		FileName: fileName,
		FileExt:  fileExt,
		FileBody: buf.Bytes(),
		FileSize: size,
		Token:    token,
	})
}

func (cli FileClient) DeleteFile(token string, fileID string) (*pb.DeleteFileRes, error) {
	return cli.service.DeleteFile(context.TODO(), &pb.DeleteFileReq{
		Token:  token,
		FileID: fileID,
	})
}
