package handler

import (
	"context"

	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/upload/config"
	"github.com/miRemid/kira/services/upload/repository"
	"github.com/pkg/errors"
)

type Handler struct {
	Repo repository.Repository
}

func (handler Handler) UploadFile(ctx context.Context, in *pb.UploadFileReq, res *pb.UploadFileRes) error {
	resp, err := handler.Repo.UploadFile(ctx, in.Owner, in.FileName, in.FileExt, in.FileSize, in.FileBody)
	if err != nil {
		res.Succ = false
		res.Msg = err.Error()
		return errors.WithMessage(err, "upload file")
	}
	res.Succ = true
	res.Msg = "upload success"
	res.File = new(pb.File)
	res.File.FileExt = resp.FileExt
	res.File.FileHash = resp.FileHash
	res.File.FileID = resp.FileID
	res.File.FileSize = resp.FileSize
	res.File.FileURL = config.Path(resp.FileID)
	res.File.FileName = resp.FileName
	return nil
}
