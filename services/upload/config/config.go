package config

import (
	"fmt"
	"log"
	"math/rand"
	"time"

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

	DOMAIN = "img.test.me"

	BUCKET_NAME = []string{"kira-1", "kira-2", "kira-3"}
)

func init() {
	pro := common.Getenv("GIN_MODE", "")
	if pro == "release" {
		DOMAIN = fmt.Sprintf("http://%s", common.Getenv("DOMAIN", "img.test.me"))
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

func Bucket() string {
	rand.Seed(time.Now().Unix())
	return BUCKET_NAME[rand.Intn(len(BUCKET_NAME))]
}

func Path(id string) string {
	return fmt.Sprintf("http://%s/image?id=%s", DOMAIN, id)
}
