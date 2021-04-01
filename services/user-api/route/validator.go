package route

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	userPattern = "^[a-zA-Z0-9_-]{4,16}$"
	pPattern    = "^.*(?=.{7,})(?=.*[0-9])(?=.*[A-Z])(?=.*[a-z])(?=.*[!@#$%^&*? ]).*$"
)

var usernameValidator validator.Func = func(fl validator.FieldLevel) bool {
	username, ok := fl.Field().Interface().(string)
	if ok {
		if valid, err := regexp.MatchString(userPattern, username); err != nil {
			return false
		} else {
			return valid
		}
	}
	return false
}

var passwordValidator validator.Func = func(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if ok {
		if len(password) < 7 {
			return false
		}
		num := `[0-9]{1}`
		a_z := `[a-z]{1}`
		A_Z := `[A-Z]{1}`
		symbol := `[!@#~$%^&*()+|_]{1}`
		if b, err := regexp.MatchString(num, password); !b || err != nil {
			return false
		}
		if b, err := regexp.MatchString(a_z, password); !b || err != nil {
			return false
		}
		if b, err := regexp.MatchString(A_Z, password); !b || err != nil {
			return false
		}
		if b, err := regexp.MatchString(symbol, password); !b || err != nil {
			return false
		}
		return true
	}

	return false
}
