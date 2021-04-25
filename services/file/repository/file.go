package repository

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"strings"

	redigo "github.com/garyburd/redigo/redis"

	"github.com/miRemid/kira/cache/redis"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/model"
	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/file/config"

	"github.com/disintegration/gift"
	"github.com/minio/minio-go/v7"
	"github.com/robfig/cron/v3"
	"github.com/segmentio/ksuid"

	"gorm.io/gorm"
)

type DeleteStruct struct {
	FileID   string `gorm:"file_id"`
	Bucket   string `gorm:"bucket"`
	FileName string `gorm:"file_name"`
	UserID   string `gorm:"-"`
	Count    int    `gorm:"-"`
}

type FileRepository interface {
	GenerateToken(ctx context.Context, userid, userName string) (string, error)
	RefreshToken(context.Context, string) (string, error)
	GetToken(context.Context, string) (string, error)
	GetHistory(context.Context, string, int64, int64) ([]*pb.UserFile, int64, error)
	DeleteFile(context.Context, string, string) error
	GetImage(ctx context.Context, in *pb.GetImageReq) ([]byte, error)
	GetDetail(ctx context.Context, fileID string) (*pb.UserFile, error)
	DeleteUser(ctx context.Context, userID string) error
	ChangeStatus(ctx context.Context, userID string, status int64) error
	CheckStatus(ctx context.Context, token string) (int64, error)
	GetUserImages(ctx context.Context, token, userID string, offset, limit int64, desc bool) ([]*pb.UserFile, int64, error)
	GetRandomFile(ctx context.Context, token string) ([]*pb.UserFile, error)
	LikeOrDislike(ctx context.Context, in *pb.FileLikeReq) (err error)
	GetLikes(ctx context.Context, userid string, offset, limit int64, desc bool) ([]*pb.UserFile, int64, error)
	Done()
}

type FileRepositoryImpl struct {
	minioCli   *minio.Client
	db         *gorm.DB
	deleteChan chan DeleteStruct
	done       chan struct{}
	c          *cron.Cron
}

func NewFileRepository(db *gorm.DB, mini *minio.Client) FileRepository {
	var res FileRepositoryImpl
	res.minioCli = mini
	res.db = db
	res.deleteChan = make(chan DeleteStruct)
	res.done = make(chan struct{}, 1)
	go res.deleteG()
	res.c = cron.New()
	res.cronInit()
	db.AutoMigrate(model.TokenUser{}, model.LikeModel{})
	return res
}

func (repo FileRepositoryImpl) token2UserID(tx *gorm.DB, token string) (string, error) {
	var userid string
	err := tx.Model(model.TokenUser{}).Select("user_id").Where("token = ?", token).First(&userid).Error
	return userid, err
}

func (repo FileRepositoryImpl) userName2UserID(tx *gorm.DB, userName string) (string, error) {
	var userid string
	err := tx.Model(model.TokenUser{}).Select("user_id").Where("user_name = ?", userName).First(&userid).Error
	return userid, err
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
			res[i].Liked = repo.findLike(conn, repo.db, res[i].FileID, userid)
		}
	}
	return res, nil
}

func (repo FileRepositoryImpl) GetUserImages(ctx context.Context, token string, userName string, offset, limit int64, desc bool) ([]*pb.UserFile, int64, error) {
	tx := repo.db.Begin()
	var notfound bool
	var userTokenID string
	if token == common.AnonyToken {
		notfound = true
	} else {
		id, err := repo.token2UserID(tx, token)
		if err != nil {
			return nil, 0, err
		}
		userTokenID = id
	}
	userID, err := repo.userName2UserID(tx, userName)
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err := tx.Model(model.FileModel{}).Where("owner = ?", userID).Count(&total).Error; err != nil {
		tx.Rollback()
		return nil, 0, err
	}
	var files = make([]*pb.UserFile, 0)
	sqlExec := `
	select ttu.user_name as user_name, tf.file_name, tf.file_width , tf.file_height ,tf.file_id ,tf.file_size ,tf.file_ext , tf.file_hash 
	from tbl_file tf left join tbl_token_user ttu 
	on tf.owner = ttu.user_id 
	where tf.owner = '%v' 
	order by tf.created_at %v LIMIT %v, %v
	`
	if desc {
		err = tx.Raw(fmt.Sprintf(sqlExec, userID, "desc", offset, limit)).Scan(&files).Error
	} else {
		err = tx.Raw(fmt.Sprintf(sqlExec, userID, "asc", offset, limit)).Scan(&files).Error
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
			files[i].Liked = repo.findLike(conn, repo.db, files[i].FileID, userTokenID)
		}
	}
	return files, total, err
}

func (repo FileRepositoryImpl) LikeOrDislike(ctx context.Context, in *pb.FileLikeReq) error {
	if in.Dislike {
		return repo.disLike(ctx, in)
	}
	return repo.like(ctx, in)
}

// Get User's Likes
func (repo FileRepositoryImpl) GetLikes(ctx context.Context, userid string, offset, limit int64, desc bool) ([]*pb.UserFile, int64, error) {
	// Get From Redis
	conn := redis.Get()
	defer conn.Close()
	var fileIDs = make([]string, 0)

	res, err := redigo.Values(conn.Do("HSCAN", common.LikeRankHash, "0", "match", userid+":*", "count", "1"))
	if err != nil {
		return nil, 0, err
	}
	values, _ := redigo.StringMap(res[1], nil)
	for k, _ := range values {
		arg := strings.Split(k, ":")
		fileIDs = append(fileIDs, arg[1])
	}
	// Get From Database
	var total int64
	if err := repo.db.Model(model.LikeModel{}).Where("user_id = ?", userid).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var ids = make([]string, 0)
	if err := repo.db.Model(model.LikeModel{}).
		Select("file_id").
		Where("user_id = ? and status = 1", userid).
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&ids).Error; err != nil {
		return nil, 0, err
	}
	fileIDs = append(fileIDs, ids...)

	// get file infomation
	// 1. get username
	var username string
	if err := repo.db.Model(model.TokenUser{}).Select("user_name").Where("user_id = ?", userid).Find(&username).Error; err != nil {
		return nil, 0, err
	}

	var files = make([]*pb.UserFile, 0)
	for _, id := range fileIDs {
		var file model.FileModel
		if err := repo.db.Model(model.FileModel{}).Where("file_id = ?", id).First(&file).Error; err != nil {
			continue
		}
		files = append(files, &pb.UserFile{
			UserName: username,
			FileName: file.FileName,
			Width:    file.FileWidth,
			Height:   file.FileHeight,
			FileID:   id,
			Liked:    true,
			Ext:      file.FileExt,
			Hash:     file.FileHash,
			FileURL:  config.Path(id),
		})
	}
	return files, total, nil
}

// Token Requests
func (repo FileRepositoryImpl) GetHistory(ctx context.Context, token string, limit, offset int64) ([]*pb.UserFile, int64, error) {
	conn := redis.Get()
	defer conn.Close()

	var total int64
	var res = make([]*pb.UserFile, 0)
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
	for i := 0; i < len(res); i++ {
		res[i].FileURL = config.Path(res[i].FileID)
		res[i].Liked = repo.findLike(conn, repo.db, res[i].FileID, owner)
	}

	tx.Commit()
	return res, total, nil
}

func (repo FileRepositoryImpl) RefreshToken(ctx context.Context, token string) (string, error) {
	log.Println("Refresh Token for: ", token)
	tx := repo.db.Begin()
	ntoken := ksuid.New().String()
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

	// 1. delete database
	if err := tx.Exec("delete from tbl_file where file_id = ?", fileID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Exec("delete from tbl_likes where file_id = ?", fileID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. delete redis
	conn := redis.Get()
	defer conn.Close()
	if _, err = conn.Do("HDEL", common.LikeRankHash, common.UserLikeFileKey(owner, fileID)); err != nil {
		tx.Rollback()
		return err
	}

	// 3. delete minio
	var bucket string
	if err := tx.Model(model.FileModel{}).Select("bucket").Where("file_id = ? and owner = ?", fileID, owner).Scan(&bucket).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := repo.minioCli.RemoveObject(ctx, bucket, fileID, minio.RemoveObjectOptions{}); err != nil {
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
	repo.c.Stop()
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
	conn := redis.Get()
	defer conn.Close()

	var offset, limit = 0, 5
	for i := 0; i < total; i += limit {
		offset += i
		var dels = make([]DeleteStruct, 0)
		if err := tx.Raw("select file_id, bucket, file_name from tbl_file where owner = ? limit ?, ?", userID, offset, limit).Scan(&dels).Error; err != nil {
			tx.Rollback()
			return err
		}
		for i := range dels {
			conn.Do("HDEL", common.LikeRankHash, common.UserLikeFileKey(userID, dels[i].FileID))
			repo.deleteChan <- dels[i]
		}
	}
	if err := tx.Exec("delete from tbl_token_user where user_id = ?", userID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Exec("delete from tbl_likes where user_id = ?", userID).Error; err != nil {
		tx.Rollback()
		return err
	}
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
	tx := repo.db.Begin()
	var token string
	if err := tx.Raw("select token from tbl_token_user where user_id = ?", userID).Scan(&token).Error; err != nil {
		log.Println(err)
		tx.Rollback()
		return "", err
	}
	return token, nil
}

func (repo FileRepositoryImpl) GetImage(ctx context.Context, in *pb.GetImageReq) ([]byte, error) {
	// 1. Get bucket
	var bucket string
	if err := repo.db.Model(model.FileModel{}).Select("bucket").Where("file_id = ?", in.FileID).Scan(&bucket).Error; err != nil {
		return nil, err
	}
	// 2. Get Files body
	obj, err := repo.minioCli.GetObject(ctx, bucket, in.FileID, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	img, _, _ := image.Decode(obj)
	g := gift.New()
	need := false
	// 3. check width and height
	if in.Width != 0 && in.Height != 0 {
		need = true
		g.Add(gift.Resize(int(in.Width), int(in.Height), gift.LanczosResampling))
	}
	// 4. check gray
	if in.Gray || in.Binary {
		need = true
		g.Add(gift.Grayscale())
	}
	if in.Binary {
		need = true
		g.Add(gift.Threshold(float32(in.Threshold)))
	}
	// 4. check blur
	if in.Blur {
		need = true
		g.Add(gift.GaussianBlur(float32(in.BlurSeed)))
	}
	var buffer bytes.Buffer
	if need {
		out := image.NewNRGBA(g.Bounds(img.Bounds()))
		g.Draw(out, img)
		err = jpeg.Encode(&buffer, out, nil)
	} else {
		err = jpeg.Encode(&buffer, img, nil)
	}
	return buffer.Bytes(), err
}

func (repo FileRepositoryImpl) GetDetail(ctx context.Context, fileID string) (*pb.UserFile, error) {
	var file = new(pb.UserFile)
	err := repo.db.Raw("select * from tbl_file where file_id = ?", fileID).Scan(&file).Error
	return file, err
}
