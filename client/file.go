package client

import (
	"context"

	"github.com/miRemid/kira/proto/pb"
	"github.com/micro/go-micro/v2/client"
)

type FileClient struct {
	service pb.FileService
}

func NewFileClient(client client.Client) *FileClient {
	var cli FileClient
	srv := pb.NewFileService("kira.micro.service.file", client)
	cli.service = srv
	return &cli
}

func (cli FileClient) GetImage(fileid, width, height string) (*pb.GetImageRes, error) {
	return cli.service.GetImage(context.TODO(), &pb.GetImageReq{
		FileID: fileid,
		Width:  width,
		Height: height,
	})
}

func (cli FileClient) GenerateToken(userid string) (*pb.TokenUserRes, error) {
	return cli.service.GenerateToken(context.TODO(), &pb.TokenUserReq{
		Userid: userid,
	})
}

func (cli FileClient) RefreshToken(token string) (*pb.TokenUserRes, error) {
	return cli.service.RefreshToken(context.TODO(), &pb.TokenReq{
		Token: token,
	})
}

func (cli FileClient) GetToken(userid string) (*pb.TokenUserRes, error) {
	return cli.service.GetToken(context.TODO(), &pb.TokenUserReq{
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

func (cli FileClient) DeleteFile(token string, fileID string) (*pb.DeleteFileRes, error) {
	return cli.service.DeleteFile(context.TODO(), &pb.DeleteFileReq{
		Token:  token,
		FileID: fileID,
	})
}

func (cli FileClient) GetDetail(fileID string) (*pb.GetDetailRes, error) {
	return cli.service.GetDetail(context.TODO(), &pb.GetDetailReq{
		FileID: fileID,
	})
}

func (client *FileClient) Ping() (*pb.Pong, error) {
	return client.service.Ping(context.TODO(), &pb.Ping{})
}
