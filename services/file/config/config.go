package config

var (
	SUPPORT_EXT = []string{"jpg", "jpeg", "png", "git"}
)

func CheckExt(ext string) bool {
	for i := 0; i < len(SUPPORT_EXT); i++ {
		if SUPPORT_EXT[i] == ext {
			return true
		}
	}
	return false
}
