package config

import (
	"fmt"
	"log"

	"github.com/miRemid/kira/common"
)

var (
	SUPPORT_EXT = []string{".jpg", ".jpeg", ".png", ".gif"}

	contentType = map[string]string{
		".jpg":  "image/jpg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
	}

	DOMAIN    = "http://api.test.me"
	IMAGE_API = DOMAIN + "/file/image/"
)

func init() {
	pro := common.Getenv("GIN_MODE", "")
	if pro == "release" {
		DOMAIN = fmt.Sprintf("http://%s", common.Getenv("DOMAIN", "api.test.me"))
		IMAGE_API = DOMAIN + "/file/image/"
	}
}

func CheckExt(ext string) bool {
	log.Print(ext)
	for i := 0; i < len(SUPPORT_EXT); i++ {
		if SUPPORT_EXT[i] == ext {
			return true
		}
	}
	return false
}

func ContentType(ext string) string {
	return contentType[ext]
}
