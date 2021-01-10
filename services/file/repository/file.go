package repository

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"time"

	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/services/file/model"
	"github.com/teris-io/shortid"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

var (
	idgen *shortid.Shortid
)

func init() {
	idgen, _ = shortid.New(8, shortid.DefaultABC, uint64(time.Now().Unix()))
}

type FileRepository interface {
	GenerateToken(context.Context, string) (string, error)
	RefreshToken(context.Context, string) (string, error)
	GetHistory(context.Context, string, int64, int64) ([]model.FileModel, int64, error)
	UploadFile(ctx context.Context, token string, fileName string, fileExt string, fileSize int64, fileBody []byte) (model.FileModel, error)
	DeleteFile(context.Context, string, string) error
	GetImage(ctx context.Context, fileID string) (model.FileModel, io.Reader, error)
}

type FileRepositoryImpl struct {
	minioCli *minio.Client
	db       *gorm.DB
}

func NewFileRepository(db *gorm.DB) (FileRepository, error) {
	var res FileRepositoryImpl
	var ssl = false
	endpoint := common.Getenv("MINIO_ENDPOINT", "127.0.0.1:9900")
	accessKey := common.Getenv("MINIO_ACCESSKEY", "kira")
	secretKey := common.Getenv("MINIO_SECRETKEY", "1234567890")
	secure := common.Getenv("MINIO_SECURE", "false")
	if secure == "true" {
		ssl = true
	}
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: ssl,
	})
	if err != nil {
		return nil, err
	}
	res.minioCli = cli
	res.db = db
	err = db.AutoMigrate(model.FileModel{}, model.TokenUser{})
	return res, err
}

// generate user's token, and create the user bucket
func (repo FileRepositoryImpl) GenerateToken(ctx context.Context, userID string) (string, error) {
	tx := repo.db.Begin()
	token := ksuid.New().String()
	log.Println("Token= ", token)
	var item model.TokenUser
	item.UserID = userID
	item.Token = token
	if err := tx.FirstOrCreate(&item).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	log.Println("Create success")
	err := repo.minioCli.MakeBucket(ctx, userID, minio.MakeBucketOptions{})
	if err != nil {
		log.Printf("Create Bucket Error %s\n", err.Error())
		tx.Rollback()
		exists, errBucketExists := repo.minioCli.BucketExists(ctx, userID)
		if errBucketExists == nil && exists {
			return "", errors.New("bucket aleardy exist")
		} else {
			return "", err
		}
	}
	tx.Commit()
	return token, nil
}

func (repo FileRepositoryImpl) RefreshToken(ctx context.Context, userID string) (string, error) {
	tx := repo.db.Begin()
	token := ksuid.New().String()
	if err := tx.Exec("update tbl_token_user set token = ? where user_id = ?", token, userID).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	return token, nil
}

func (repo FileRepositoryImpl) GetHistory(ctx context.Context, token string, limit, offset int64) ([]model.FileModel, int64, error) {
	var total int64
	var res = make([]model.FileModel, 0)
	var tx = repo.db.Begin()
	// 1. get user_id
	var userid string
	if err := tx.Raw("select user_id from tbl_token_user where token = ?", token).Scan(&userid).Error; err != nil {
		tx.Rollback()
		return res, total, err
	}
	// 2. get count
	if err := tx.Raw("select COUNT(*) from tbl_file group by user_id = ?", userid).Scan(&total).Error; err != nil {
		tx.Rollback()
		return res, total, err
	}
	// 3. get files list
	if err := tx.Raw("select * from tbl_file where user_id = ? limit ?, ?", userid, offset, offset+limit).Scan(&res).Error; err != nil {
		tx.Rollback()
		return res, total, err
	}
	tx.Commit()
	return res, total, nil
}

func (repo FileRepositoryImpl) UploadFile(ctx context.Context,
	token, fileName, fileExt string,
	fileSize int64, fileBody []byte) (model.FileModel, error) {
	log.Print(len(fileBody))
	var res model.FileModel

	var tx = repo.db.Begin()
	// 1. get user_id
	var userid string
	if err := tx.Raw("select user_id from tbl_token_user where token = ?", token).Scan(&userid).Error; err != nil {
		tx.Rollback()
		return res, err
	}
	res.UserID = userid
	// 2. generate file's sha1 hash
	hash := sha1.New()
	reader := bytes.NewReader(fileBody)
	if _, err := io.Copy(hash, reader); err != nil {
		tx.Rollback()
		return res, err
	}
	hashInBytes := hash.Sum(nil)[:20]
	res.FileHash = hex.EncodeToString(hashInBytes)
	res.FileName = fileName
	res.FileSize = fileSize
	res.FileExt = fileExt
	res.FileID, _ = idgen.Generate()

	reader.Seek(0, 0)
	// 3. upload into minio
	_, err := repo.minioCli.PutObject(ctx, userid, fileName, reader, int64(fileSize), minio.PutObjectOptions{})
	if err != nil {
		tx.Rollback()
		return res, err
	}

	// 4. insert record into database
	if err := tx.Model(model.FileModel{}).FirstOrCreate(&res).Error; err != nil {
		tx.Rollback()
		return res, err
	}
	tx.Commit()
	return res, nil
}

func (repo FileRepositoryImpl) DeleteFile(ctx context.Context, token string, fileID string) error {
	var tx = repo.db.Begin()
	// 1. get user_id
	var userid string
	if err := tx.Raw("select user_id from tbl_token_user where token = ?", token).Scan(&userid).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 2. get file's name
	var fileName string
	if err := tx.Raw("select file_name from tbl_file where file_id = ?", fileID).Scan(&fileName).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 2. delete minio file
	if err := repo.minioCli.RemoveObject(ctx, userid, fileName, minio.RemoveObjectOptions{}); err != nil {
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
	obj, err := repo.minioCli.GetObject(ctx, file.UserID, file.FileName, minio.GetObjectOptions{})
	if err != nil {
		tx.Rollback()
		return file, nil, err
	}
	return file, obj, nil
}
