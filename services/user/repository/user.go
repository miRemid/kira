package repository

import (
	"fmt"
	"log"

	"github.com/miRemid/kira/model"
	authClient "github.com/miRemid/kira/services/auth/client"
	fileClient "github.com/miRemid/kira/services/file/client"
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
	Refresh(userid string) (string, error)
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
	pwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(pwd)
	user.UserID = xid.New().String()

	res, err := repo.fileCli.GenerateToken(user.UserID)
	if err != nil || !res.Succ {
		return err
	}

	user.Token = res.Token
	user.Role = "normal"

	if err := repo.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (repo UserRepositoryImpl) Signin(username, password string) (string, error) {
	tx := repo.db.Begin()

	// 1. get user model
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

func (repo UserRepositoryImpl) UserInfo(userid string) (model.UserModel, error) {
	var user model.UserModel
	tx := repo.db.Begin()
	// Get User Info
	if err := tx.Raw("select * from tbl_user where user_id = ?", userid).Scan(&user).Error; err != nil {
		tx.Rollback()
		return user, err
	}
	// Get User Token
	if res, err := repo.fileCli.GetToken(userid); err != nil {
		tx.Rollback()
		return user, err
	} else {
		user.Token = res.Token
	}
	tx.Commit()
	return user, nil
}
