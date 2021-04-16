package handler

import (
	"context"
	"log"

	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/file/config"
	"github.com/miRemid/kira/services/file/repository"
	"github.com/pkg/errors"
)

type FileServiceHandler struct {
	Repo repository.FileRepository
}

func (handler FileServiceHandler) GetHotLikeRank(ctx context.Context, in *pb.Empty, res *pb.HotLikeRankList) error {
	files, err := handler.Repo.GetHotLikeRank(ctx)
	res.Files = files
	return err
}

func (handler FileServiceHandler) LikeOrDislike(ctx context.Context, in *pb.FileLikeReq, res *pb.Response) error {
	err := handler.Repo.LikeOrDislike(ctx, in.Userid, in.Fileid, in.Dislike)
	if err == nil {
		res.Code = int64(response.StatusOK)
		res.Message = "successful"
	} else if err == response.ErrAlreadyLike {
		res.Code = int64(response.StatusAlreadyLike)
		res.Message = err.Error()
	} else {
		res.Code = int64(response.StatusInternalError)
		res.Message = err.Error()
	}
	return nil
}

func (handler FileServiceHandler) GetRandomFile(ctx context.Context, in *pb.Empty, res *pb.RandomFiles) error {
	resp, err := handler.Repo.GetRandomFile(ctx)
	if err != nil {
		return err
	}
	for i := 0; i < len(resp); i++ {
		resp[i].FileURL = config.Path(resp[i].FileID)
	}
	res.Files = resp
	return nil
}

func (handler FileServiceHandler) GetUserImages(ctx context.Context, in *pb.GetUserImagesReq, res *pb.GetUserImagesRes) error {
	items, total, err := handler.Repo.GetUserImages(ctx, in.Userid, in.Offset, in.Limit, in.Desc)
	if err != nil {
		return errors.WithMessage(err, "get user images")
	}
	var files = make([]*pb.File, 0)
	for i := 0; i < len(items); i = i + 1 {
		var file pb.File
		file.FileExt = items[i].FileExt
		file.FileName = items[i].FileName
		file.FileID = items[i].FileID
		file.FileHash = items[i].FileHash
		file.FileWidth = items[i].FileWidth
		file.FileHeight = items[i].FileHeight
		file.FileSize = int64(items[i].FileSize)
		file.FileURL = config.Path(items[i].FileID)
		files = append(files, &file)
	}
	res.Files = files
	res.Total = total
	return nil
}

func (handler FileServiceHandler) GetImage(ctx context.Context, in *pb.GetImageReq, res *pb.GetImageRes) error {
	data, err := handler.Repo.GetImage(ctx, in.FileID, in.Width, in.Height)
	if err != nil {
		res.Msg = err.Error()
		res.Succ = false
		return errors.WithMessage(err, "get image")
	}

	res.Msg = "get success"
	res.Succ = true
	res.Image = data
	return nil
}

func (handler FileServiceHandler) GenerateToken(ctx context.Context, in *pb.TokenUserReq, res *pb.TokenUserRes) error {
	log.Printf("Generate Token for %s\n", in.Userid)
	token, err := handler.Repo.GenerateToken(ctx, in.Userid, in.UserName)
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
func (handler FileServiceHandler) RefreshToken(ctx context.Context, in *pb.TokenReq, res *pb.TokenUserRes) error {
	token, err := handler.Repo.RefreshToken(ctx, in.Token)
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
func (handler FileServiceHandler) GetToken(ctx context.Context, in *pb.TokenUserReq, res *pb.TokenUserRes) error {
	token, err := handler.Repo.GetToken(ctx, in.Userid)
	if err != nil {
		res.Msg = err.Error()
		res.Succ = false
		return errors.WithMessage(err, "get token")
	}
	res.Msg = "get token success"
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
		file.FileWidth = items[i].FileWidth
		file.FileHeight = items[i].FileHeight
		file.FileSize = int64(items[i].FileSize)
		file.FileURL = config.Path(items[i].FileID)
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

func (handler FileServiceHandler) GetDetail(ctx context.Context, in *pb.GetDetailReq, res *pb.GetDetailRes) error {
	resp, err := handler.Repo.GetDetail(ctx, in.FileID)
	if err != nil {
		res.Succ = false
		res.Msg = err.Error()
		return errors.WithMessage(err, "get detail")
	}
	res.Succ = true
	res.Msg = "get success"
	res.File = new(pb.File)
	res.File.FileID = resp.FileID
	res.File.FileURL = config.Path(resp.FileID)
	res.File.FileHash = resp.FileHash
	res.File.FileName = resp.FileName
	res.File.FileExt = resp.FileExt
	res.File.FileSize = resp.FileSize
	res.File.FileWidth = resp.FileWidth
	res.File.FileHeight = resp.FileHeight
	return nil
}

func (handler FileServiceHandler) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) error {
	return handler.Repo.DeleteUser(ctx, in.UserID)
}

func (handler FileServiceHandler) DeleteAnony(ctx context.Context, in *pb.DeleteFileReq) error {
	log.Println("Delete Anony File, FileID = ", in.FileID)
	return handler.Repo.DeleteFile(ctx, "", in.FileID)
}

func (handler FileServiceHandler) Ping(ctx context.Context, in *pb.Ping, res *pb.Pong) error {
	res.Code = 0
	res.Name = "file"
	res.Message = "ok"
	return nil
}

func (handler FileServiceHandler) ChangeTokenStatus(ctx context.Context, in *pb.ChangeTokenStatusReq, res *pb.ChangeTokenStatusRes) error {
	if err := handler.Repo.ChangeStatus(ctx, in.Userid, in.Status); err != nil {
		res.Succ = false
		res.Msg = err.Error()
		return err
	}
	res.Succ = true
	res.Msg = "change successful"
	return nil
}
func (handler FileServiceHandler) CheckTokenStatus(ctx context.Context, in *pb.CheckTokenStatusReq, res *pb.CheckTokenStatusRes) error {
	status, err := handler.Repo.CheckStatus(ctx, in.Token)
	if err != nil {
		res.Succ = false
		res.Msg = err.Error()
		return err
	}
	res.Succ = true
	res.Msg = "check successful"
	res.Status = status
	return nil
}
