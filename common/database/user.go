package database

import (
	"github.com/miRemid/kira/model"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func InitAdmin(username, password string, db *gorm.DB) {
	var user model.UserModel
	if err := db.Raw("select * from tbl_user where user_name = ?", username).Scan(&user).Error; err == gorm.ErrRecordNotFound {
		user.UserName = username
		pwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user.Password = string(pwd)
		user.UserID = xid.New().String()
		user.Role = "admin"
		db.Create(&user)
	}
}
