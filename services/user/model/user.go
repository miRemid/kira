package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	UserID   string
	UserName string
	Password string
	Role     string
	Token    string
}

func (UserModel) TableName() string {
	return "tbl_user"
}

func (user UserModel) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
}
