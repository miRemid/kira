package repository

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/miRemid/kira/cache/redis"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/model"
	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/upload/config"
	"github.com/micro/go-micro/v2"
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
		mini.MakeBucket(ctx, name, minio.MakeBucketOptions{})
	}
	mini.MakeBucket(ctx, common.AnonyBucket, minio.MakeBucketOptions{})
}

type Repository interface {
	UploadFile(ctx context.Context, owner string, fileName string, fileExt string, fileSize int64, fileWidth, fileHeight string, anony bool, fileBody []byte) (model.FileModel, error)
}

type RepositoryImpl struct {
	mini *minio.Client
	db   *gorm.DB
	pub  micro.Event
}

func NewRepository(db *gorm.DB, mini *minio.Client, pub micro.Event) Repository {
	bucket(mini)
	db.AutoMigrate(model.FileModel{})
	repo := RepositoryImpl{
		mini: mini,
		db:   db,
		pub:  pub,
	}
	go repo.deleteAnony()
	return repo
}

func (repo RepositoryImpl) deleteAnony() {
	// for every 1 hour, check the redis
	log.Println("Start Delete Anony Files Time Tricker: 1 hour")
	tricker := time.NewTicker(time.Hour)
	defer tricker.Stop()
	for t := range tricker.C {
		log.Println("Tricker Delete Anony Files...")
		log.Println("Start get fileid from redis")
		// get the redis
		conn := redis.Get()
		timestamp := t.Unix()
		res, err := redigo.StringMap(conn.Do("ZRANGEBYSCORE", common.AnonymousKey, "-inf", timestamp, "withscores"))
		if err != nil {
			log.Println("Get fileids from redis err: ", err)
		} else {
			log.Println("Get fileid from redis successful, length = ", len(res))
			for fileid := range res {
				log.Printf("Send %s to the nats message queue", fileid)
				// insert into the nats
				repo.pub.Publish(context.TODO(), &pb.DeleteFileReq{
					Token:  "",
					FileID: fileid,
				})
				// remove from the redis
				conn.Do("ZREM", common.AnonymousKey, fileid)
			}
		}
		conn.Close()
	}
}

func (repo RepositoryImpl) UploadFile(ctx context.Context,
	token, fileName, fileExt string,
	fileSize int64, fileWidth, fileHeight string, anony bool, fileBody []byte) (model.FileModel, error) {
	var res model.FileModel

	id, _ := idgen.Generate()

	// 2. generate file's sha1 hash
	hash := sha1.New()
	reader := bytes.NewReader(fileBody)
	if _, err := io.Copy(hash, reader); err != nil {
		return res, err
	}

	bucket := config.Bucket(anony)
	var tx = repo.db.Begin()
	var userid string = common.AnonyBucket
	if !anony {
		tx.Model(model.TokenUser{}).Select("user_id").Where("token = ?", token).First(&userid)
	}
	res.Owner = userid
	hashInBytes := hash.Sum(nil)[:20]
	res.FileWidth = fileWidth
	res.FileHeight = fileHeight
	res.FileHash = hex.EncodeToString(hashInBytes)
	res.FileName = fileName
	res.FileSize = fileSize
	res.FileExt = fileExt
	res.FileID = id
	res.Bucket = bucket
	// 5. if anony upload, insert into redis delay queue
	if anony {
		res.Anony = true
		conn := redis.Get()
		defer conn.Close()
		// save 5 day for the anony upload
		delay := time.Now().Add(time.Hour * 24 * 5).Unix()
		conn.Do("ZADD", common.AnonymousKey, delay, id)
	}
	// 4. insert record into database
	if err := tx.Model(model.FileModel{}).Create(&res).Error; err != nil {
		tx.Rollback()
		return res, err
	}
	reader.Seek(0, 0)
	// 3. upload into minio
	_, err := repo.mini.PutObject(ctx, bucket, id, reader, int64(fileSize), minio.PutObjectOptions{})
	if err != nil {
		tx.Rollback()
		return res, err
	}
	tx.Commit()
	return res, nil
}
