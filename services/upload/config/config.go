package config

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
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

	DOMAIN = "test.me:5000"

	BUCKET_NAME = []string{"kira-1", "kira-2", "kira-3"}

	TEMP_DIR = "tmp"
)

func init() {
	pro := common.Getenv("GIN_MODE", "")
	if pro == gin.ReleaseMode {
		DOMAIN = common.Getenv("DOMAIN", "img.test.me")
	}
	TEMP_DIR = filepath.Join(common.Getenv("TEMP_DIR", "./"), TEMP_DIR)
	rand.Seed(time.Now().UnixNano())
}

func CheckExt(ext string) bool {
	_, ok := contentType[ext]
	return ok
}

func ContentType(ext string) string {
	return contentType[ext]
}

func Bucket(anony bool) string {
	if anony {
		return common.AnonyBucket
	}
	index := rand.Intn(len(BUCKET_NAME))
	return BUCKET_NAME[index]
}

func Path(id string) string {
	return fmt.Sprintf("http://%s/image?id=%s", DOMAIN, id)
}
