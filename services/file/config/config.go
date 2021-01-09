package config

import "log"

var (
	SUPPORT_EXT = []string{".jpg", ".jpeg", ".png", ".gif"}

	contentType = map[string]string{
		".jpg":  "image/jpg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
	}
)

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
