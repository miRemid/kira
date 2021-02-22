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
	log.Println("Init Admin Account: ", username)
	user.UserName = username
	pwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(pwd)
	user.UserID = xid.New().String()
	user.Role = "admin"
	if err := db.Create(&user).Error; err != nil {
		log.Println(err)
	}
}
