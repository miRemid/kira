package handler

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/miRemid/kira/services/file/pb"
	"github.com/miRemid/kira/services/file/repository"
	"github.com/pkg/errors"
)

type FileServiceHandler struct {
	Repo repository.FileRepository
}

func (handler FileServiceHandler) GetImage(ctx context.Context, in *pb.GetImageReq, res *pb.GetImageRes) error {
	file, reader, err := handler.Repo.GetImage(ctx, in.FileID)
	if err != nil {
		res.Msg = err.Error()
		res.Succ = false
		return errors.WithMessage(err, "get image")
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil {
		res.Msg = err.Error()
		res.Succ = false
		return errors.WithMessage(err, "get image")
	}
	res.Msg = "get success"
	res.Succ = true
	res.Image = buf.Bytes()
	res.FileExt = file.FileExt
	res.FileName = file.FileName
	return nil
}

func (handler FileServiceHandler) GenerateToken(ctx context.Context, in *pb.TokenUserReq, res *pb.TokenUserRes) error {
	log.Printf("Generate Token for %s\n", in.Userid)
	token, err := handler.Repo.GenerateToken(ctx, in.Userid)
	if err != nil {
		log.Printf("Generate Token Err: %v\n", err)
		res.Msg = err.Error()
		res.Succ = false
		return errors.WithMessage(err, "generate token")
	}
	log.Printf("Generate Token %s\n", token)
	res.Msg = "generate success"
	res.Succ = true
	res.Token = token
	return nil
}
func (handler FileServiceHandler) RefreshToken(ctx context.Context, in *pb.TokenUserReq, res *pb.TokenUserRes) error {
	token, err := handler.Repo.RefreshToken(ctx, in.Userid)
	if err != nil {
		res.Msg = err.Error()
		res.Succ = false
		return errors.WithMessage(err, "refresh token")
	}
	res.Msg = "refresh success"
	res.Succ = true
	res.Token = token
	return nil
}
func (handler FileServiceHandler) GetHistory(ctx context.Context, in *pb.GetHistoryReq, res *pb.GetHistoryRes) error {
	items, total, err := handler.Repo.GetHistory(ctx, in.Token, in.Limit, in.Offset)
	if err != nil {
		res.Succ = false
		res.Msg = err.Error()
		return errors.WithMessage(err, "get histroy")
	}
	var files = make([]*pb.File, 0)
	for i := 0; i < len(items); i = i + 1 {
		var file pb.File
		file.FileExt = items[i].FileExt
		file.FileName = items[i].FileName
		file.FileID = items[i].FileID
		file.FileHash = items[i].FileHash
		file.FileSize = int64(items[i].FileSize)
		files = append(files, &file)
	}
	res.Total = total
	res.Files = files
	res.Msg = "get success"
	res.Succ = true
	return nil
}
func (handler FileServiceHandler) DeleteFile(ctx context.Context, in *pb.DeleteFileReq, res *pb.DeleteFileRes) error {
	err := handler.Repo.DeleteFile(ctx, in.Token, in.FileID)
	if err != nil {
		res.Succ = false
		res.Msg = err.Error()
		return errors.WithMessage(err, "delete file")
	}
	res.Succ = true
	res.Msg = "delete success"
	return nil
}
func (handler FileServiceHandler) UploadFile(ctx context.Context, in *pb.UploadFileReq, res *pb.UploadFileRes) error {
	log.Println(len(in.FileBody))
	fileInfo, err := handler.Repo.UploadFile(ctx, in.Token, in.FileName, in.FileExt, in.FileSize, in.FileBody)
	if err != nil {
		log.Println(err)
		res.Succ = false
		res.Msg = err.Error()
		return errors.WithMessage(err, "upload file")
	}
	res.Succ = true
	res.Msg = "upload success"
	var file = new(pb.File)
	file.FileExt = fileInfo.FileExt
	file.FileHash = fileInfo.FileHash
	file.FileID = fileInfo.FileID
	file.FileName = fileInfo.FileName
	file.FileSize = fileInfo.FileSize
	res.File = file
	return nil
}
