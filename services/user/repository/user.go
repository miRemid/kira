package repository

import (
	authClient "github.com/miRemid/kira/services/auth/client"
	fileClient "github.com/miRemid/kira/services/file/client"
	"github.com/miRemid/kira/services/user/model"
	"github.com/micro/go-micro/v2/client"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	Signup(username, password string) error
	Signin(username, password string) (string, error)
	UserInfo(userid string) (model.UserModel, error)
}

type UserRepositoryImpl struct {
	db      *gorm.DB
	authCli *authClient.AuthClient
	fileCli *fileClient.FileClient
}

func NewUserRepository(service client.Client, db *gorm.DB) (UserRepository, error) {
	ac := authClient.NewAuthClient(service)
	fc := fileClient.NewFileClient(service)
	err := db.AutoMigrate(model.UserModel{})
	return UserRepositoryImpl{
		db:      db,
		authCli: ac,
		fileCli: fc,
	}, err
}

func (repo UserRepositoryImpl) Signup(username, password string) error {
	var user model.UserModel
	user.UserName = username
	pwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(pwd)
	user.UserID = xid.New().String()

	res, err := repo.fileCli.GenerateToken(user.UserID)
	if err != nil || !res.Succ {
		return err
	}
	user.Token = res.Token
	user.Role = "normal"

	if err := repo.db.FirstOrCreate(&user).Error; err != nil {
		return err
	}

	return nil
}

func (repo UserRepositoryImpl) Signin(username, password string) (string, error) {
	tx := repo.db.Begin()

	// 1. get user model
	var user model.UserModel
	if err := tx.Raw("select * from tbl_user where user_name = ?", username).Scan(&user).Error; err != nil {
		tx.Rollback()
		return "", err
	}

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

func (repo UserRepositoryImpl) UserInfo(userid string) (model.UserModel, error) {
	var user model.UserModel
	tx := repo.db.Begin()
	if err := tx.Raw("select * from tbl_user where user_id = ?", userid).Scan(&user).Error; err != nil {
		tx.Rollback()
		return user, err
	}
	tx.Commit()
	return user, nil
}
