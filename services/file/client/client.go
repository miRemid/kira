package client

import (
	"bytes"
	"context"
	"io"
	"log"
	"mime/multipart"
	"time"

	"github.com/miRemid/kira/services/file/pb"
	microClient "github.com/micro/go-micro/v2/client"
)

type FileClient struct {
	service pb.FileService
}

func NewFileClient(service microClient.Client) *FileClient {
	var cli FileClient
	srv := pb.NewFileService("kira.micro.service.file", service)
	cli.service = srv
	return &cli
}

func (cli FileClient) GetImage(fileid string) (*pb.GetImageRes, error) {
	return cli.service.GetImage(context.TODO(), &pb.GetImageReq{
		FileID: fileid,
	})
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

func (cli FileClient) UploadFile(token string, fileName, fileExt string, file multipart.File) (*pb.UploadFileRes, error) {
	var buf bytes.Buffer
	size, _ := io.Copy(&buf, file)
	log.Println(len(buf.Bytes()))
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancel()
	return cli.service.UploadFile(ctx, &pb.UploadFileReq{
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

func (cli FileClient) GetDetail(fileID string) (*pb.GetDetailRes, error) {
	return cli.service.GetDetail(context.TODO(), &pb.GetDetailReq{
		FileID: fileID,
	})
}
