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

func (j *JWT) CreateToken(userName string, expireUnix int64, customClaims ...CustomClaimOption) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = userName
	claims["exp"] = expireUnix

	for _, cc := range customClaims {
		cc.apply(claims)
	}

	t, err := token.SignedString([]byte(j.key))
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

// ParseAuditPlanName used by echo middleware which only verify api request to audit plan related.
func ParseAuditPlanName(tokenString string) (string, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	}
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	apn, ok := claims["apn"]
	if !ok {
		return "", jwt.NewValidationError("unknown token", jwt.ValidationErrorClaimsInvalid)
	}
	return apn.(string), nil
}
