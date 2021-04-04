package config

import (
	"fmt"
	"path/filepath"

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
}

func CheckExt(ext string) bool {
	_, ok := contentType[ext]
	return ok
}

func ContentType(ext string) string {
	return contentType[ext]
}

func Bucket() string {
	return "kira-1"
}

func Path(id string) string {
	return fmt.Sprintf("http://%s/image?id=%s", DOMAIN, id)
}
