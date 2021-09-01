package utils

import (
	"github.com/dgrijalva/jwt-go"
)

// TODO: Using configuration to set jwt secret
const JWTSecret = "secret"

type JWT struct {
	key []byte
}

func NewJWT(signingKey []byte) *JWT {
	return &JWT{
		key: signingKey,
	}
}

func (j *JWT) CreateToken(userName string, expireUnix int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = userName
	claims["exp"] = expireUnix

	t, err := token.SignedString([]byte(j.key))
	if err != nil {
		return "", err
	}
	return t, nil
}
