package config

import (
	"fmt"

	"github.com/miRemid/kira/common"
)

var (
	contentType = map[string]string{
		".jpg":  "image/jpg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
	}

	DOMAIN = "test.me:5000"
)

func init() {
	pro := common.Getenv("GIN_MODE", "")
	if pro == "release" {
		DOMAIN = common.Getenv("DOMAIN", "img.test.me")
	}
}

func ContentType(ext string) string {
	return contentType[ext]
}

func Path(id string) string {
	return fmt.Sprintf("http://%s/image?id=%s", DOMAIN, id)
}
