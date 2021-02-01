package repository

import (
	"context"
	"io"
	"time"

	"github.com/miRemid/kira/model"
	"github.com/teris-io/shortid"

	"github.com/minio/minio-go/v7"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

var (
	idgen *shortid.Shortid
)

func init() {
	idgen, _ = shortid.New(8, shortid.DefaultABC, uint64(time.Now().Unix()))
}

type DeleteStruct struct {
	FileID   string `gorm:"file_id"`
	Bucket   string `gorm:"bucket"`
	FileName string `gorm:"file_name"`
	Count    int    `gorm:"-"`
}

type FileRepository interface {
	GenerateToken(context.Context, string) (string, error)
	RefreshToken(context.Context, string) (string, error)
	GetToken(context.Context, string) (string, error)
	GetHistory(context.Context, string, int64, int64) ([]model.FileModel, int64, error)
	DeleteFile(context.Context, string, string) error
	GetImage(ctx context.Context, fileID string) (model.FileModel, io.Reader, error)
	GetDetail(ctx context.Context, fileID string) (model.FileModel, error)
	DeleteUser(ctx context.Context, userID string) error
	Done()
}

type FileRepositoryImpl struct {
	minioCli   *minio.Client
	db         *gorm.DB
	deleteChan chan DeleteStruct
	done       chan struct{}
}

func NewFileRepository(db *gorm.DB, mini *minio.Client) FileRepository {
	var res FileRepositoryImpl
	res.minioCli = mini
	res.db = db
	res.deleteChan = make(chan DeleteStruct)
	res.done = make(chan struct{}, 1)
	go res.deleteG()
	db.AutoMigrate(model.TokenUser{})
	return res
}

func (repo FileRepositoryImpl) deleteG() {
	for {
		select {
		case item := <-repo.deleteChan:
			repo.minioCli.RemoveObject(context.Background(), item.Bucket, item.FileName+"-"+item.FileID, minio.RemoveObjectOptions{})
			repo.db.Raw("delete from tbl_file where file_id = ?", item.FileID)
			break
		case <-repo.done:
			return
		default:
			break
		}
	}
}

func (repo FileRepositoryImpl) Done() {
	repo.done <- struct{}{}
	close(repo.done)
	close(repo.deleteChan)
}

func (repo FileRepositoryImpl) DeleteUser(ctx context.Context, userID string) error {
	defer func() {
		if r := recover(); r != nil {

		}
	}()
	// 1. Get User FileID
	var total int
	tx := repo.db.Begin()
	if err := tx.Raw("select COUNT(*) from tbl_file where owner = ?", userID).Scan(&total).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 2. 批量删除
	var offset, limit = 0, 5
	for i := 0; i < total; i += limit {
		offset += i
		var dels = make([]DeleteStruct, 0)
		if err := tx.Raw("select file_id, bucket, file_name from tbl_file where user_id = ? limit ?, ?", userID, offset, limit).Scan(&dels).Error; err != nil {
			tx.Rollback()
			return err
		}
		for i := range dels {
			repo.deleteChan <- dels[i]
		}
	}
	tx.Raw("delete from tbl_token_user where user_id = ?", userID)
	tx.Commit()
	return nil
}

// generate user's token, and create the user bucket
func (repo FileRepositoryImpl) GenerateToken(ctx context.Context, userID string) (string, error) {
	tx := repo.db.Begin()
	token := ksuid.New().String()
	var item model.TokenUser
	item.UserID = userID
	item.Token = token
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	return token, nil
}

func (repo FileRepositoryImpl) RefreshToken(ctx context.Context, token string) (string, error) {
	tx := repo.db.Begin()
	ntoken := ksuid.New().String()
	if err := tx.Exec("update tbl_token_user set token = ? where token = ?", ntoken, token).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	return token, nil
}

func (repo FileRepositoryImpl) GetToken(ctx context.Context, userID string) (string, error) {
	tx := repo.db.Begin()
	var token string
	if err := tx.Raw("select token from tbl_token_user where user_id = ?", userID).Scan(&token).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	return token, nil
}

func (repo FileRepositoryImpl) GetHistory(ctx context.Context, owner string, limit, offset int64) ([]model.FileModel, int64, error) {
	var total int64
	var res = make([]model.FileModel, 0)
	var tx = repo.db.Begin()
	if err := tx.Raw("select COUNT(*) from tbl_file where owner = ?", owner).Scan(&total).Error; err != nil {
		tx.Rollback()
		return res, total, err
	}
	// 3. get files list
	if err := tx.Raw("select * from tbl_file where owner = ? limit ?, ?", owner, offset, limit).Scan(&res).Error; err != nil {
		tx.Rollback()
		return res, total, err
	}
	tx.Commit()
	return res, total, nil
}

func (repo FileRepositoryImpl) DeleteFile(ctx context.Context, owner string, fileID string) error {
	var tx = repo.db.Begin()
	// 1. get bucket
	var bucket, fileName string
	row := tx.Raw("select bucket, file_name from tbl_file where file_id = ? and owner = ?", fileID, owner).Row()
	if err := row.Scan(&bucket, &fileName); err != nil {
		tx.Rollback()
		return err
	}
	// 2. delete minio file
	if err := repo.minioCli.RemoveObject(ctx, bucket, fileName+"-"+fileID, minio.RemoveObjectOptions{}); err != nil {
		tx.Rollback()
		return err
	}
	// 3. delete database record
	if err := tx.Exec("delete from tbl_file where file_id = ?", fileID).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (repo FileRepositoryImpl) GetImage(ctx context.Context, fileID string) (model.FileModel, io.Reader, error) {
	var file model.FileModel
	tx := repo.db.Begin()
	if err := tx.Raw("select * from tbl_file where file_id = ?", fileID).Scan(&file).Error; err != nil {
		tx.Rollback()
		return file, nil, err
	}
	// 2. Get Files body
	obj, err := repo.minioCli.GetObject(ctx, file.Bucket, file.FileName+"-"+file.FileID, minio.GetObjectOptions{})
	if err != nil {
		tx.Rollback()
		return file, nil, err
	}
	return file, obj, nil
}

func (repo FileRepositoryImpl) GetDetail(ctx context.Context, fileID string) (model.FileModel, error) {
	var file model.FileModel
	err := repo.db.Raw("select * from tbl_file where file_id = ?", fileID).Scan(&file).Error
	return file, err
}
