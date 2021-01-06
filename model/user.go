package model

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	UserID       string
	UserName     string
	UserPassword string
	Role         string
	FileToken    string
}
