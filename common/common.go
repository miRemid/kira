package common

import "os"

func Getenv(key, replace string) string {
	res := os.Getenv(key)
	if res == "" {
		return replace
	}
	return res
}
