package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/miRemid/kira/cache/redis"
	"github.com/miRemid/kira/client"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/common/database"
	"github.com/miRemid/kira/model"
	"github.com/miRemid/kira/proto/pb"
	"github.com/micro/go-micro/v2"
	mClient "github.com/micro/go-micro/v2/client"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Signup(username, password string) error
	Signin(username, password string) (string, error)
	UserInfo(username string) (model.UserModel, error)

	GetUserList(ctx context.Context, limit, offset int64) ([]model.UserModel, int64, error)
	ChangeUserStatus(ctx context.Context, userid string, status int64) error
	DeleteUser(ctx context.Context, userid string) error
	GetUserImages(ctx context.Context, userName string, offset, limit int64, desc bool) (*pb.GetUserImagesRes, error)

	ChangePassword(ctx context.Context, userid, old, raw string) error
	LikeOrDislike(ctx context.Context, req *pb.FileLikeReq) error
	GetUserToken(ctx context.Context, userid string) (string, error)
}

type UserRepositoryImpl struct {
	db      *gorm.DB
	authCli *client.AuthClient
	fileCli *client.FileClient

	pub micro.Event
}

func NewUserRepository(service mClient.Client, db *gorm.DB, pub micro.Event) (UserRepository, error) {
	ac := client.NewAuthClient(service)
	fc := client.NewFileClient(service)
	err := db.AutoMigrate(model.UserModel{})
	database.InitAdmin(common.Getenv("ADMIN_USERNAME", "miosuki"), common.Getenv("ADMIN_PASSWORD", "QAZplm%123"), db)
	return UserRepositoryImpl{
		db:      db,
		authCli: ac,
		fileCli: fc,
		pub:     pub,
	}, err
}

func (repo UserRepositoryImpl) GetUserToken(ctx context.Context, userid string) (string, error) {
	res, err := repo.fileCli.GetToken(userid)
	return res.Token, err
}

func (repo UserRepositoryImpl) LikeOrDislike(ctx context.Context, req *pb.FileLikeReq) error {
	// 1. call file service to incr count for fileid
	file, err := repo.fileCli.Service.LikeOrDislike(ctx, req)
	if err != nil {
		return err
	}
	// 2. insert file's infomation into the user's like hash set
	// key = user_id_likes
	key := common.UserLikeKey(req.Userid)
	conn := redis.Get()
	defer conn.Close()
	buffer, _ := json.Marshal(file)
	_, err = conn.Do("SET", key, buffer)
	return err
}

func (repo UserRepositoryImpl) GetUserImages(ctx context.Context, userName string, offset, limit int64, desc bool) (*pb.GetUserImagesRes, error) {
	// 1. get user id
	var userid string
	if err := repo.db.Model(model.UserModel{}).Select("user_id").Where("user_name = ?", userName).First(&userid).Error; err != nil {
		log.Println("Get User ", userName, " failed: ", err)
		return nil, err
	}
	// 2. rpc call
	return repo.fileCli.GetUserImages(&pb.GetUserImagesReq{
		Userid: userid,
		Offset: offset,
		Limit:  limit,
		Desc:   desc,
	})
}

func (repo UserRepositoryImpl) ChangePassword(ctx context.Context, userid, old, npwd string) error {
	log.Println("Change Password for userid = ", userid)
	tx := repo.db.Begin()
	var user model.UserModel
	if err := tx.Model(user).Where("user_id = ?", userid).First(&user).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return errors.New("username not found")
		}
		return errors.WithMessage(err, "get user model")
	}
	log.Printf("UserID = %s, UserName = %s", user.UserID, user.UserName)
	if !user.CheckPassword(old) {
		tx.Rollback()
		return errors.New("old password incorrect")
	}
	user.Password = user.GeneratePassword(npwd)
	if err := tx.Model(user).Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	log.Println("Change Password for userid = ", userid, " successful")
	return nil
}

func (repo UserRepositoryImpl) Refresh(userid string) (string, error) {
	res, err := repo.fileCli.RefreshToken(userid)
	if err != nil || !res.Succ {
		return "", errors.New("User Service: Refresh failed")
	}
	return res.Token, nil
}

func (repo UserRepositoryImpl) Signup(username, password string) error {
	var user model.UserModel
	user.UserName = username
	if err := repo.db.Model(user).Where("user_name = ?", username).First(&user).Error; err == nil {
		return fmt.Errorf("username '%s' already exists", username)
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	user.Password = string(user.GeneratePassword(password))
	user.UserID = xid.New().String()

	res, err := repo.fileCli.GenerateToken(user.UserID, username)
	if err != nil || !res.Succ {
		return err
	}

	user.Token = res.Token
	user.Role = "normal"
	user.Status = 1

	if err := repo.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (repo UserRepositoryImpl) Signin(username, password string) (string, error) {
	tx := repo.db.Begin()

	// get user model
	var user model.UserModel
	if err := tx.Model(model.UserModel{}).Where("user_name = ?", username).First(&user).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return "", errors.New("username not found")
		}
		return "", errors.WithMessage(err, "get user model")
	}

	log.Printf("UserID = %s, UserName = %s", user.UserID, username)

	// 2. check password
	if !user.CheckPassword(password) {
		tx.Rollback()
		return "", errors.New("wrong password")
	}

	// 3. generate token
	res, err := repo.authCli.Auth(user.UserID, user.Role)
	if err != nil || !res.Succ {
		tx.Rollback()
		return "", errors.WithMessage(err, res.Msg)
	}
	tx.Commit()
	return res.Token, nil
}

func (repo UserRepositoryImpl) UserInfo(username string) (model.UserModel, error) {
	log.Println("receive", username)
	var user model.UserModel
	tx := repo.db.Begin()
	// Get User Info
	if err := tx.Table(user.TableName()).Where("user_name = ?", username).First(&user).Error; err != nil {
		tx.Rollback()
		return user, err
	}
	tx.Commit()
	return user, nil
}

func (repo UserRepositoryImpl) GetUserList(ctx context.Context, limit, offset int64) ([]model.UserModel, int64, error) {
	var total int64
	var res = make([]model.UserModel, 0)
	tx := repo.db.Begin()
	if err := tx.Raw(`select COUNT(*) from tbl_user where role = "normal" `).Scan(&total).Error; err != nil {
		tx.Rollback()
		return res, total, err
	}
	if err := tx.Raw(`select * from tbl_user where role = "normal" limit ?, ?`, offset, limit).Scan(&res).Error; err != nil {
		tx.Rollback()
		return res, total, err
	}
	tx.Commit()
	return res, total, nil
}

func (repo UserRepositoryImpl) DeleteUser(ctx context.Context, userid string) error {
	var role string
	if err := repo.db.Raw("select role from tbl_user where user_id = ?", userid).Scan(&role).Error; err != nil {
		return err
	}
	if role == "admin" {
		return errors.New("Cannot delete admin user")
	}
	if err := repo.db.Exec("delete from tbl_user where user_id = ?", userid).Error; err != nil {
		return err
	}
	// publish delete event
	return repo.pub.Publish(ctx, &pb.DeleteUserRequest{
		UserID: userid,
	})
}

// Change user status, allow user to signin, get history item, delete item, but cannot upload
func (repo UserRepositoryImpl) ChangeUserStatus(ctx context.Context, userid string, status int64) error {
	log.Println("Change status for userid = ", userid)
	_, err := repo.fileCli.ChangeStatus(userid, status)
	if err != nil {
		return err
	}
	return repo.db.Model(&model.UserModel{}).Where("user_id = ?", userid).Update("status", status).Error
}
