package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	UserID   string `gorm:"column:user_id;index:idx_user_id,unique"`
	UserName string `gorm:"column:user_name;index:idx_user_name,unique"`
	Password string `gorm:"column:password"`
	Role     string `gorm:"column:role"`
	Status   int64  `gorm:"column:status"` // 0 inactive, 1 active
	Token    string `gorm:"-"`
}

func (UserModel) TableName() string {
	return "tbl_user"
}

func (user UserModel) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
}

func (user UserModel) GeneratePassword(password string) string {
	pwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(pwd)
}
