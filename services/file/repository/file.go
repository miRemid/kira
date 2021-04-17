package repository

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"log"
	"strconv"
	"time"

	redigo "github.com/garyburd/redigo/redis"

	"github.com/miRemid/kira/cache/redis"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/response"
	"github.com/miRemid/kira/model"
	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/file/config"

	"github.com/minio/minio-go/v7"
	"github.com/nfnt/resize"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type DeleteStruct struct {
	FileID   string `gorm:"file_id"`
	Bucket   string `gorm:"bucket"`
	FileName string `gorm:"file_name"`
	Count    int    `gorm:"-"`
}

type FileRepository interface {
	GenerateToken(ctx context.Context, userid, userName string) (string, error)
	RefreshToken(context.Context, string) (string, error)
	GetToken(context.Context, string) (string, error)
	GetHistory(context.Context, string, int64, int64) ([]model.FileModel, int64, error)
	DeleteFile(context.Context, string, string) error
	GetImage(ctx context.Context, fileID, width, height string) ([]byte, error)
	GetDetail(ctx context.Context, fileID string) (model.FileModel, error)
	DeleteUser(ctx context.Context, userID string) error
	ChangeStatus(ctx context.Context, userID string, status int64) error
	CheckStatus(ctx context.Context, token string) (int64, error)
	GetUserImages(ctx context.Context, token, userID string, offset, limit int64, desc bool) ([]*pb.UserFile, int64, error)
	GetRandomFile(ctx context.Context, token string) ([]*pb.UserFile, error)
	LikeOrDislike(ctx context.Context, userid string, fileid string, dislike bool) (err error)
	GetHotLikeRank(ctx context.Context, token string) ([]*pb.UserFile, error)
	GetLikes(ctx context.Context, userid string, offset, limit int64, desc bool) ([]*pb.UserFile, int64, error)
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

func (repo FileRepositoryImpl) token2UserID(tx *gorm.DB, token string) (string, error) {
	var userid string
	err := tx.Model(model.TokenUser{}).Select("user_id").Where("token = ?", token).First(&userid).Error
	return userid, err
}

// Normal And Anony Requests
func (repo FileRepositoryImpl) GetHotLikeRank(ctx context.Context, token string) ([]*pb.UserFile, error) {
	var userid string
	var notfound bool
	if token == common.AnonyToken {
		notfound = true
	} else {
		userid, _ = repo.token2UserID(repo.db, token)
	}

	conn := redis.Get()
	defer conn.Close()

	// 1. get top 10 from the set
	fileScore, err := redigo.Strings(conn.Do("ZREVRANGE", common.LikeRankKey, 0, 9, "WITHSCORES"))
	if err != nil {
		return nil, err
	}
	var files = make([]*pb.UserFile, 0)
	for i := 0; i < len(fileScore)/2; i = i + 2 {
		id, score := fileScore[i], fileScore[i+1]
		var file = new(pb.UserFile)
		if err = repo.db.Raw(`
		select ttu.user_name, tf.file_name, tf.file_id, tf.file_width, tf.file_height
		from tbl_file tf left join tbl_token_user ttu on tf.owner = ttu.user_id
		where tf.file_id = ?`, id).Scan(file).Error; err != nil {
			return nil, err
		}
		file.FileURL = config.Path(file.FileID)
		s, _ := strconv.Atoi(score)
		file.Likes = int64(s)
		files = append(files, file)
		// 2. if token != ""
		if !notfound {
			userKey := common.UserLikeKey(userid)
			exist, err := redigo.Bool(conn.Do("HEXISTS", userKey, files[i].FileID))
			if err != nil || !exist {
				files[i].Liked = false
			} else {
				files[i].Liked = true
			}
		}
	}
	return files, nil
}

func (repo FileRepositoryImpl) GetRandomFile(ctx context.Context, token string) ([]*pb.UserFile, error) {
	var userid string
	var notfound bool
	if token == common.AnonyToken {
		notfound = true
	} else {
		userid, _ = repo.token2UserID(repo.db, token)
	}

	var res = make([]*pb.UserFile, 0)
	if err := repo.db.Raw(`select ttu.user_name, tf.file_name, tf.file_id, tf.file_width, tf.file_height, tf.anony 
	from tbl_file tf left join tbl_token_user ttu on tf.owner = ttu.user_id 
	where tf.anony != 1 
	and tf.id >= ((SELECT MAX(tf2.id) from tbl_file tf2) - (select MIN(tf3.id) from tbl_file tf3)) * RAND() + (select MIN(tu.id) from tbl_user tu)  
	limit 20`).Scan(&res).Error; err != nil {
		return nil, err
	}
	conn := redis.Get()
	defer conn.Close()
	for i := 0; i < len(res); i++ {
		res[i].FileURL = config.Path(res[i].FileID)
		// 1. get rank
		index, err := redigo.Int64(conn.Do("ZRANK", common.LikeRankKey, res[i].FileID))
		if err != nil {
			res[i].Likes = 0
		} else {
			likesMap, err := redigo.Int64Map(conn.Do("ZRANGE", common.LikeRankKey, index, index, "WITHSCORES"))
			if err != nil {
				res[i].Likes = 0
			} else {
				res[i].Likes = likesMap[res[i].FileID]
			}
		}
		// 2. if token != ""
		if !notfound {
			exist, err := getUserFileLikeStatus(conn, res[i].FileID, userid)
			if err != nil || !exist {
				res[i].Liked = false
			} else {
				res[i].Liked = true
			}
		}
	}
	return res, nil
}

func (repo FileRepositoryImpl) GetUserImages(ctx context.Context, token string, userID string, offset, limit int64, desc bool) ([]*pb.UserFile, int64, error) {
	var notfound bool
	if token == common.AnonyToken {
		notfound = true
	}
	log.Println(userID, offset, limit, desc)
	tx := repo.db.Begin()
	var total int64
	if err := tx.Model(model.FileModel{}).Where("owner = ?", userID).Count(&total).Error; err != nil {
		tx.Rollback()
		return nil, 0, err
	}
	var files = make([]*pb.UserFile, 0)
	var err error
	if desc {
		err = tx.Model(model.FileModel{}).Where("owner = ?", userID).Order("created_at desc").Limit(int(limit)).Offset(int(offset)).Find(&files).Error
	} else {
		err = tx.Model(model.FileModel{}).Where("owner= ?", userID).Order("created_at asc").Limit(int(limit)).Offset(int(offset)).Find(&files).Error
	}
	if err != nil {
		tx.Rollback()
	}
	conn := redis.Get()
	defer conn.Close()
	for i := 0; i < len(files); i++ {
		files[i].FileURL = config.Path(files[i].FileID)
		// 1. get rank
		index, err := redigo.Int64(conn.Do("ZRANK", common.LikeRankKey, files[i].FileID))
		if err != nil {
			files[i].Likes = 0
		} else {
			likesMap, err := redigo.Int64Map(conn.Do("ZRANGE", common.LikeRankKey, index, index, "WITHSCORES"))
			if err != nil {
				files[i].Likes = 0
			} else {
				files[i].Likes = likesMap[files[i].FileID]
			}
		}
		// 2. if token != ""
		if !notfound {
			exist, err := getUserFileLikeStatus(conn, files[i].FileID, userID)
			if err != nil || !exist {
				files[i].Liked = false
			} else {
				files[i].Liked = true
			}
		}
	}
	return files, total, err
}

// Signin User Requests
func (repo FileRepositoryImpl) LikeOrDislike(ctx context.Context, userid string, fileid string, dislike bool) (err error) {
	conn := redis.Get()
	defer conn.Close()
	key := common.UserLikeKey(userid)
	exist, err := getUserFileLikeStatus(conn, fileid, userid)
	if err != nil {
		return err
	}
	if !dislike && exist {
		return response.ErrAlreadyLike
	} else if dislike && !exist {
		return response.ErrNotLike
	}

	var offset = 1
	if dislike {
		offset = -1
	}
	var file = new(pb.UserFile)
	if err = repo.db.Raw(`
	select ttu.user_name, tf.file_name, tf.file_id, tf.file_width, tf.file_height
	from tbl_file tf left join tbl_token_user ttu on tf.owner = ttu.user_id
	where tf.file_id = ?`, fileid).Scan(file).Error; err != nil {
		return err
	}
	// 1. incr zset rank list
	if _, err = conn.Do("zincrby", common.LikeRankKey, offset, fileid); err != nil {
		return err
	}
	if dislike {
		// remove the id from the user's set
		conn.Do("ZREM", key, fileid)
	} else {
		// save file's infomation
		setFileIntoHash(conn, file)
		// 2. insert userid_like hash table, key = fileid, value = file
		datetime := time.Now().Unix()
		conn.Do("ZADD", key, datetime, file.FileID)
	}
	return err
}

// Get User's Likes
func (repo FileRepositoryImpl) GetLikes(ctx context.Context, userid string, offset, limit int64, desc bool) ([]*pb.UserFile, int64, error) {
	conn := redis.Get()
	defer conn.Close()

	userKey := common.UserLikeKey(userid)

	total, err := redigo.Int64(conn.Do("ZCARD", userKey))
	if err != nil {
		return nil, 0, err
	}

	cmd := "ZRANGE"
	if desc {
		cmd = "ZREVRANGE"
	}
	log.Printf("UserID = %s, Offset = %v, Limit = %v", userKey, offset, offset+limit)
	fileids, err := redigo.Strings(conn.Do(cmd, userKey, offset, offset+limit))
	if err != nil {
		return nil, 0, err
	}
	var files = make([]*pb.UserFile, 0)
	log.Println(len(fileids))
	for i := 0; i < len(fileids); i++ {
		files = append(files, getFileFromHash(conn, userid, fileids[i]))
	}
	return files, total, nil
}

// Token Requests
func (repo FileRepositoryImpl) GetHistory(ctx context.Context, token string, limit, offset int64) ([]model.FileModel, int64, error) {
	conn := redis.Get()
	defer conn.Close()

	var total int64
	var res = make([]model.FileModel, 0)
	var tx = repo.db.Begin()
	owner, err := repo.token2UserID(tx, token)
	if err != nil {
		tx.Rollback()
		return res, total, err
	}
	if err := tx.Raw("select COUNT(*) from tbl_file where owner = ?", owner).Scan(&total).Error; err != nil {
		tx.Rollback()
		return res, total, err
	}
	// 3. get files list
	if err := tx.Raw("select * from tbl_file where owner = ? order by created_at, id desc limit ?, ? ", owner, offset, limit).Scan(&res).Error; err != nil {
		tx.Rollback()
		return res, total, err
	}
	log.Println(owner)
	for i := 0; i < len(res); i++ {
		exist, err := getUserFileLikeStatus(conn, res[i].FileID, owner)
		if err != nil || !exist {
			res[i].Liked = false
		} else {
			res[i].Liked = true
		}
		log.Println(res[i].FileID, res[i].Liked)
	}

	tx.Commit()
	return res, total, nil
}

func (repo FileRepositoryImpl) RefreshToken(ctx context.Context, token string) (string, error) {
	log.Println("Refresh Token for: ", token)
	tx := repo.db.Begin()
	ntoken := ksuid.New().String()
	userid, err := repo.token2UserID(tx, token)
	if err != nil {
		log.Println("Refresh Token, Get infomation err: ", err)
		tx.Rollback()
		return "", err
	}
	conn := redis.Get()
	if _, err := conn.Do("DEL", userid); err != nil {
		log.Println("Delete key ", userid, " failed: ", err)
		tx.Rollback()
		return "", err
	}
	if err := tx.Exec("update tbl_token_user set token = ? where token = ?", ntoken, token).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	log.Println("Refresh Token: ", ntoken)
	return ntoken, nil
}

func (repo FileRepositoryImpl) DeleteFile(ctx context.Context, token string, fileID string) error {
	var tx = repo.db.Begin()
	owner, err := repo.token2UserID(tx, token)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 1. get bucket
	var bucket string
	if err := tx.Model(model.FileModel{}).Select("bucket").Where("file_id = ? and owner = ?", fileID, owner).Scan(&bucket).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 2. delete minio file
	if err := repo.minioCli.RemoveObject(ctx, bucket, fileID, minio.RemoveObjectOptions{}); err != nil {
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

// Admin Requests
func (repo FileRepositoryImpl) ChangeStatus(ctx context.Context, userID string, status int64) error {
	// modify userID's token status
	// 0: suspend, 1: active
	log.Println("Change status for userid = ", userID)
	return repo.db.Model(&model.TokenUser{}).Where("user_id = ?", userID).Update("status", status).Error
}

func (repo FileRepositoryImpl) CheckStatus(ctx context.Context, token string) (int64, error) {
	var status int64
	err := repo.db.Model(&model.TokenUser{}).Select("status").Where("token = ?", token).Scan(&status).Error
	return status, err
}

func (repo FileRepositoryImpl) deleteG() {
	for {
		select {
		case item, ok := <-repo.deleteChan:
			if !ok {
				return
			}
			repo.minioCli.RemoveObject(context.Background(), item.Bucket, item.FileID, minio.RemoveObjectOptions{})
			repo.db.Exec("delete from tbl_file where file_id = ?", item.FileID)
		case <-repo.done:
			return
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
			log.Println(r)
		}
	}()
	log.Printf("Rcv Message From Nats: userid=%v", userID)
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
		if err := tx.Raw("select file_id, bucket, file_name from tbl_file where owner = ? limit ?, ?", userID, offset, limit).Scan(&dels).Error; err != nil {
			tx.Rollback()
			return err
		}
		for i := range dels {
			repo.deleteChan <- dels[i]
		}
	}
	if err := tx.Exec("delete from tbl_token_user where user_id = ?", userID).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	log.Println("Delete userid: ", userID)
	return nil
}

// RPC
// generate user's token, and create the user bucket
func (repo FileRepositoryImpl) GenerateToken(ctx context.Context, userID, userName string) (string, error) {
	log.Println("Generate Token For UserID: ", userID)
	tx := repo.db.Begin()
	token := ksuid.New().String()
	var item model.TokenUser
	item.UserID = userID
	item.Token = token
	item.UserName = userName
	item.Status = 1
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	log.Printf("UserID: %v, Token: %v\n", userID, token)
	return token, nil
}

// This method is for user service
func (repo FileRepositoryImpl) GetToken(ctx context.Context, userID string) (string, error) {
	log.Println("Get Token For UserID: ", userID)
	tx := repo.db.Begin()
	var token string
	if err := tx.Raw("select token from tbl_token_user where user_id = ?", userID).Scan(&token).Error; err != nil {
		log.Println(err)
		tx.Rollback()
		return "", err
	}
	log.Println("UserID: ", userID, "; Token: ", token)
	return token, nil
}

func (repo FileRepositoryImpl) GetImage(ctx context.Context, fileID, width, height string) ([]byte, error) {
	// 1. Get bucket
	var bucket string
	if err := repo.db.Model(model.FileModel{}).Select("bucket").Where("file_id = ?", fileID).Scan(&bucket).Error; err != nil {
		return nil, err
	}
	// 2. Get Files body
	obj, err := repo.minioCli.GetObject(ctx, bucket, fileID, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	img, _, _ := image.Decode(obj)
	w, err := strconv.Atoi(width)
	if err != nil {
		return nil, err
	}
	h, err := strconv.Atoi(height)
	if err != nil {
		return nil, err
	}
	var out image.Image = img
	if w != 0 && h != 0 {
		out = resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	}

	var buffer bytes.Buffer
	err = jpeg.Encode(&buffer, out, nil)
	return buffer.Bytes(), err
}

func (repo FileRepositoryImpl) GetDetail(ctx context.Context, fileID string) (model.FileModel, error) {
	var file model.FileModel
	err := repo.db.Raw("select * from tbl_file where file_id = ?", fileID).Scan(&file).Error
	return file, err
}
