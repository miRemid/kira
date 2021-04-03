package database

import (
	"log"

	"github.com/miRemid/kira/model"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func InitAdmin(username, password string, db *gorm.DB) {
	var user model.UserModel
	pwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// 1. get user
	err := db.Model(user).Where("user_name = ?", username).First(&user).Error
	// if not found, create
	if err == gorm.ErrRecordNotFound {
		log.Println("Init Admin Account: ", username)
		user.UserName = username
		user.Password = string(pwd)
		user.UserID = xid.New().String()
		user.Role = "admin"
		user.Status = 1
		if err := db.Create(&user).Error; err != nil {
			log.Println(err)
		}
		// if found, update password
	} else {
		user.Password = string(pwd)
		if err = db.Model(user).Save(&user).Error; err != nil {
			log.Println(err)
		}
	}
}
