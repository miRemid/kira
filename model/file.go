package model

import "gorm.io/gorm"

type FileModel struct {
	gorm.Model
	FileID   string
	FileName string
	FileExt  string
	FileSize int64
	FileHash string
}
