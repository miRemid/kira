package model

import (
	"fmt"
	"hash/adler32"

	"gorm.io/gorm"
)

type FileModel struct {
	gorm.Model
	Owner      string `gorm:"column:owner;index:idx_owner"`
	FileID     string `gorm:"column:file_id;index:idx_file_id"`
	FileName   string `gorm:"column:file_name;"`
	FileExt    string `gorm:"column:file_ext;"`
	FileSize   int64  `gorm:"column:file_size;"`
	FileHash   string `gorm:"column:file_hash;index:idx_file_hash"`
	FileURL    string `gorm:"-"`
	Likes      int64  `gorm:"-"`
	Liked      bool   `gorm:"-"`
	Bucket     string `gorm:"column:bucket"`
	Anony      bool   `gorm:"column:anony"`
	FileWidth  string `gorm:"column:file_width"`
	FileHeight string `gorm:"column:file_height"`
}

func (FileModel) TableName() string {
	return "tbl_file"
}

func (file *FileModel) HashID() {
	h := adler32.New()
	file.FileID = fmt.Sprintf("%x", h.Sum([]byte(file.FileName)))
}

type TokenUser struct {
	gorm.Model
	UserID   string `gorm:"column:user_id;index:idx_user_id,unique"`
	Token    string `gorm:"column:token;index:idx_token"`
	UserName string `gorm:"user_name;index:idx_name"`
	Status   int    `gorm:"column:status"`
}

func (TokenUser) TableName() string {
	return "tbl_token_user"
}
