package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

var JWTSecretKey = []byte("secret")

func setJWTSecretKey(key []byte) {
	JWTSecretKey = key
}

type JWT struct {
	key []byte
}

func NewJWT(signingKey []byte) *JWT {
	return &JWT{
		key: signingKey,
	}
}

func (j *JWT) CreateToken(userName string, expireUnix int64, customClaims ...CustomClaimOption) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	// claims can only be jwt.MapClaims
	//nolint:forcetypeassert
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = userName
	claims["exp"] = expireUnix

	for _, cc := range customClaims {
		cc.apply(claims)
	}

	t, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}
	return t, nil
}

type CustomClaimOption interface {
	apply(jwt.MapClaims)
}

type funcCustomClaimOption struct {
	f func(jwt.MapClaims)
}

func newFuncCustomClaimOption(f func(jwt.MapClaims)) *funcCustomClaimOption {
	return &funcCustomClaimOption{f: f}
}

func (fdo *funcCustomClaimOption) apply(do jwt.MapClaims) {
	fdo.f(do)
}

func WithAuditPlanName(name string) CustomClaimOption {
	return newFuncCustomClaimOption(func(mc jwt.MapClaims) {
		mc["apn"] = name
	})
}

// ParseAuditPlanToken used by echo middleware which only verify api request to audit plan related.
func ParseAuditPlanToken(tokenString string) (string, string, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return JWTSecretKey, nil
	}
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok {
			if e.Errors != jwt.ValidationErrorExpired {
				return "", "", err
			}
		}
	}
	// claims can only be jwt.MapClaims
	//nolint:forcetypeassert
	claims := token.Claims.(jwt.MapClaims)
	apn, ok := claims["apn"]
	if !ok {
		return "", "", jwt.NewValidationError("unknown token", jwt.ValidationErrorClaimsInvalid)
	}
	userName, ok := claims["name"]
	if !ok {
		return "", "", jwt.NewValidationError("unknown token", jwt.ValidationErrorClaimsInvalid)
	}
	return apn.(string), userName.(string), nil
}

func GetUserNameFromJWTToken(token string) (string, error) {
	type MyClaims struct {
		Username string `json:"name"`
		jwt.StandardClaims
	}

	t, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("parse token failed: %v", err)
	}

	if claims, ok := t.Claims.(*MyClaims); ok && t.Valid {
		return claims.Username, nil
	}
	return "", fmt.Errorf("can not get out user name")
}
