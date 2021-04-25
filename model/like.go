package model

import "gorm.io/gorm"

type LikeModel struct {
	gorm.Model
	FileID string `gorm:"index;column:file_id"`
	UserID string `gorm:"index;column:user_id"`
	Status int64  `gorm:"column:status"`
}

func (LikeModel) TableName() string {
	return "tbl_likes"
}
