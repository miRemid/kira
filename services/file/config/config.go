package config

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/miRemid/kira/common"
)

var (
	contentType = map[string]string{
		".jpg":  "image/jpg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
	}

	DOMAIN = "api.test.me"
)

func init() {
	pro := common.Getenv("GIN_MODE", "")
	if pro == gin.ReleaseMode {
		DOMAIN = common.Getenv("DOMAIN", "img.test.me")
	}
}

func ContentType(ext string) string {
	return contentType[ext]
}

func Path(id string) string {
	return fmt.Sprintf("http://%s/image?id=%s", DOMAIN, id)
}
