package repository

import (
	"context"
	"log"
	"strings"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/miRemid/kira/cache/redis"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/model"
	"github.com/miRemid/kira/proto/pb"
	"gorm.io/gorm"
)

// Like service
func (repo FileRepositoryImpl) like(ctx context.Context, in *pb.FileLikeReq) error {
	conn := redis.Get()
	defer conn.Close()

	key := common.UserLikeFileKey(in.Userid, in.Fileid)
	// 1. check redis or database
	exist, _ := redigo.Bool(conn.Do("HEXISTS", common.LikeRankHash, key))
	if exist {
		return response.ErrAlreadyLike
	}
	var status int64
	if err := repo.db.Model(model.LikeModel{}).
		Select("status").
		Where("file_id = ? and user_id = ?", in.Fileid, in.Userid).
		First(&status).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}
	if status == 1 {
		return response.ErrAlreadyLike
	}
	// 2. insert into redis
	conn.Do("HSET", common.LikeRankHash, key, "1")
	return nil
}

func (repo FileRepositoryImpl) disLike(ctx context.Context, in *pb.FileLikeReq) error {
	conn := redis.Get()
	defer conn.Close()
	key := common.UserLikeFileKey(in.Userid, in.Fileid)
	// 1. check redis
	exist, _ := redigo.Bool(conn.Do("HEXISTS", common.LikeRankHash, key))
	if exist {
		// remove
		_, err := conn.Do("HDEL", common.LikeRankHash, key)
		return err
	}
	// 2. check database
	tx := repo.db.Begin()
	var status int64
	if err := tx.Model(model.LikeModel{}).
		Select("status").
		Where("file_id = ? and user_id = ?", in.Fileid, in.Userid).
		First(&status).Error; err != nil {
		tx.Rollback()
		return err
	}
	if status == 1 {
		if err := tx.Exec("delete from tbl_likes where file_id = ? and user_id = ?", in.Fileid, in.Userid).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (repo FileRepositoryImpl) findLike(conn redigo.Conn, db *gorm.DB, fileid, userid string) bool {
	key := common.UserLikeFileKey(userid, fileid)
	// 1. check redis
	exist, err := redigo.Bool(conn.Do("HEXISTS", common.LikeRankHash, key))
	if err != nil {
		return false
	}
	if exist {
		return true
	}
	// 2. check database
	var status int64
	if err := db.Model(model.LikeModel{}).Select("status").Where("file_id = ? and user_id = ?", fileid, userid).First(&status).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Println(err)
		}
		return false
	}
	return status == 1
}

func (repo FileRepositoryImpl) cronInit() {
	log.Println("Init cron tasks...")
	repo.c.AddFunc("@hourly", repo.saveRedisToDatabase)
	repo.c.Start()
	log.Println("Init cron tasks finished...")
}

func (repo FileRepositoryImpl) saveRedisToDatabase() {
	log.Println("Save redis' likes to database trigger")
	// 1. get all keys from redis
	// 2. insert into the database
	conn := redis.Get()
	defer conn.Close()

	tx := repo.db.Begin()
	values, _ := redigo.Strings(conn.Do("HKEYS", common.LikeRankHash))
	for _, uf := range values {
		argv := strings.Split(uf, ":")
		tx.Model(model.LikeModel{}).Create(&model.LikeModel{
			FileID: argv[1],
			UserID: argv[0],
			Status: 1,
		})
	}
	conn.Do("DEL", common.LikeRankHash)
	tx.Commit()
}
