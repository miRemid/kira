package token

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	// ErrTokenMalformed token
	ErrTokenMalformed = errors.New("token malfromed")
	// ErrTokenExpired token
	ErrTokenExpired = errors.New("token expired")
	// ErrTokenNotValidYet token
	ErrTokenNotValidYet = errors.New("token not valid")
	// ErrTokenInvalid token
	ErrTokenInvalid = errors.New("token invalid")
)

type AuthClaims struct {
	jwt.StandardClaims
	UserID string
	Role   string
}

type AuthControl struct {
	screct string
}

func NewAuthControl(screct string) *AuthControl {
	return &AuthControl{
		screct: screct,
	}
}

func (control *AuthControl) GenerateToken(claims *AuthClaims) *jwt.Token {
	jwt.TimeFunc = time.Now
	claims.StandardClaims.ExpiresAt = time.Now().Add(time.Hour * time.Duration(24*7)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token
}

func (control *AuthControl) ValidToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(control.screct), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}
	return token, nil
}

func (control *AuthControl) Refresh(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(control.screct), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return control.GenerateToken(claims), nil
	}
	return nil, ErrTokenInvalid
}

func (control AuthControl) GetScrect() string {
	return control.screct
}
