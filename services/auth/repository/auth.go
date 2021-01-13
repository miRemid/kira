package repository

import (
	"context"
	"crypto/rsa"

	"github.com/dgrijalva/jwt-go"
	"github.com/miRemid/kira/services/auth/token"
)

type AuthRepository interface {
	Auth(ctx context.Context, userid, username string) (string, error)
	Valid(ctx context.Context, token string) (*jwt.Token, error)
	Refresh(ctx context.Context, token string) (string, error)
}

type AuthRepositoryImpl struct {
	center *token.AuthControl
}

func (auth AuthRepositoryImpl) Auth(ctx context.Context, userid, role string) (string, error) {
	claims := token.AuthClaims{
		UserID: userid,
		Role:   role,
	}
	token := auth.center.GenerateToken(&claims)
	return token.SignedString(auth.center.GetPri())
}

func (auth AuthRepositoryImpl) Valid(ctx context.Context, tokenString string) (*jwt.Token, error) {
	return auth.center.ValidToken(tokenString)
}

func (auth AuthRepositoryImpl) Refresh(ctx context.Context, tokenString string) (string, error) {
	token, err := auth.center.Refresh(tokenString)
	if err != nil {
		return tokenString, err
	}
	return token.SignedString(auth.center.GetPri())
}

func NewAuthRepositoryImpl(pubKey *rsa.PublicKey, priKey *rsa.PrivateKey) AuthRepository {
	return AuthRepositoryImpl{
		center: token.NewAuthControl(pubKey, priKey),
	}
}
