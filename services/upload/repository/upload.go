package repository

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"time"

	"github.com/miRemid/kira/model"
	"github.com/miRemid/kira/services/upload/config"
	"github.com/minio/minio-go/v7"
	"github.com/teris-io/shortid"
	"gorm.io/gorm"
)

var (
	idgen *shortid.Shortid
)

func init() {
	idgen, _ = shortid.New(8, shortid.DefaultABC, uint64(time.Now().Unix()))
}

func bucket(mini *minio.Client) {
	ctx := context.Background()
	for _, name := range config.BUCKET_NAME {
		exists, errExists := mini.BucketExists(ctx, name)
		if errExists == nil && exists {
			continue
		}
		mini.MakeBucket(context.TODO(), name, minio.MakeBucketOptions{})
	}
}

type Repository interface {
	UploadFile(ctx context.Context, owner string, fileName string, fileExt string, fileSize int64, fileBody []byte) (model.FileModel, error)
}

type RepositoryImpl struct {
	mini *minio.Client
	db   *gorm.DB
}

func NewRepository(db *gorm.DB, mini *minio.Client) Repository {
	bucket(mini)
	db.AutoMigrate(model.FileModel{})
	return RepositoryImpl{
		mini: mini,
		db:   db,
	}
}

func (repo RepositoryImpl) UploadFile(ctx context.Context,
	owner, fileName, fileExt string,
	fileSize int64, fileBody []byte) (model.FileModel, error) {
	var res model.FileModel
	var tx = repo.db.Begin()
	id, _ := idgen.Generate()

	// 2. generate file's sha1 hash
	hash := sha1.New()
	reader := bytes.NewReader(fileBody)
	if _, err := io.Copy(hash, reader); err != nil {
		tx.Rollback()
		return res, err
	}
	bucket := config.Bucket()
	reader.Seek(0, 0)
	// 3. upload into minio
	_, err := repo.mini.PutObject(ctx, bucket, id, reader, int64(fileSize), minio.PutObjectOptions{})
	if err != nil {
		tx.Rollback()
		return res, err
	}
	res.Owner = owner
	hashInBytes := hash.Sum(nil)[:20]
	res.FileHash = hex.EncodeToString(hashInBytes)
	res.FileName = fileName
	res.FileSize = fileSize
	res.FileExt = fileExt
	res.FileID = id
	// 4. insert record into database
	if err := tx.Model(model.FileModel{}).Create(&res).Error; err != nil {
		tx.Rollback()
		return res, err
	}
	tx.Commit()
	return res, nil
}
