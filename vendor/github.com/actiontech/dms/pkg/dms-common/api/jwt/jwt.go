package jwt

import (
	"fmt"
	"strconv"
	"time"

	dmsCommonV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"

	jwtOld "github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v4"
)

type EchoContextGetter interface {
	// Get retrieves data from the context.
	Get(key string) interface{}
}

type CustomClaimFunc func(claims jwt.MapClaims)

const (
	JWTUserId        = "uid"
	JWTUsername      = "name"
	JWTExpiredTime   = "exp"
	JWTAuditPlanName = "apn"
)

func GenJwtToken(customClaims ...CustomClaimFunc) (tokenStr string, err error) {
	var mapClaims = jwt.MapClaims{
		"iss":          "actiontech dms",
		JWTExpiredTime: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	for _, claimFunc := range customClaims {
		claimFunc(mapClaims)
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)

	// Generate encoded token and send it as response.
	tokenStr, err = token.SignedString(dmsCommonV1.JwtSigningKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign the token: %v", err)
	}
	return tokenStr, nil
}

func WithUserId(userId string) CustomClaimFunc {
	return func(claims jwt.MapClaims) {
		claims[JWTUserId] = userId
	}
}

func WithUserName(name string) CustomClaimFunc {
	return func(claims jwt.MapClaims) {
		claims[JWTUsername] = name
	}
}

func WithAuditPlanName(name string) CustomClaimFunc {
	return func(claims jwt.MapClaims) {
		claims[JWTAuditPlanName] = name
	}
}

func WithExpiredTime(duration time.Duration) CustomClaimFunc {
	return func(claims jwt.MapClaims) {
		claims[JWTExpiredTime] = jwt.NewNumericDate(time.Now().Add(duration))
	}
}

func ParseUidFromJwtTokenStr(tokenStr string) (uid string, err error) {
	token, err := parseJwtTokenStr(tokenStr)
	if err != nil {
		return "", err
	}

	userId, err := ParseUserUidStrFromToken(token)
	if err != nil {
		return "", fmt.Errorf("get user id from token failed, err: %v", err)
	}

	return userId, nil
}

func parseJwtTokenStr(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if signMethod256, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		} else if signMethod256 != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}

		return dmsCommonV1.JwtSigningKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token failed: %v", err)
	}

	return token, nil
}

// ParseAuditPlanName used by echo middleware which only verify api request to audit plan related.
func ParseAuditPlanName(tokenStr string) (string, error) {
	token, err := parseJwtTokenStr(tokenStr)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to convert token claims to jwt")
	}

	auditPlanName, ok := claims[JWTAuditPlanName]
	if !ok {
		return "", jwt.NewValidationError("unknown token", jwt.ValidationErrorClaimsInvalid)
	}

	return fmt.Sprintf("%v", auditPlanName), nil
}

func GetUserFromContext(c EchoContextGetter) (uid int64, err error) {
	if c.Get("user") == nil {
		return 0, fmt.Errorf("user not found in context")
	}

	// Gets user token from the context.
	u, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return 0, fmt.Errorf("failed to convert user from jwt token")
	}
	return ParseUserFromToken(u)
}

func GetUserUidStrFromContext(c EchoContextGetter) (uid string, err error) {
	if c.Get("user") == nil {
		return "", fmt.Errorf("user not found in context")
	}

	// Gets user token from the context.
	u, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return "", fmt.Errorf("failed to convert user from jwt token")
	}
	return ParseUserUidStrFromToken(u)
}

// 由于sqle的go版本为1.16，无法使用github.com/golang-jwt/jwt/v4，本方法为sqle兼容
func GetUserUidStrFromContextWithOldJwt(c EchoContextGetter) (uid string, err error) {
	if c.Get("user") == nil {
		return "", fmt.Errorf("user not found in context")
	}

	// Gets user token from the context.
	u, ok := c.Get("user").(*jwtOld.Token)
	if !ok {
		return "", fmt.Errorf("failed to convert user from jwt token")
	}
	return ParseUserUidStrFromTokenWithOldJwt(u)
}

func ParseUserFromToken(token *jwt.Token) (uid int64, err error) {
	uidStr, err := ParseUserUidStrFromToken(token)
	if err != nil {
		return 0, err
	}
	uid, err = strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse user id: %v", err)
	}
	return uid, nil
}

func ParseUserUidStrFromToken(token *jwt.Token) (uid string, err error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to convert token claims to jwt")
	}

	uidStr := fmt.Sprintf("%v", claims[JWTUserId])
	if uidStr == "" {
		return "", fmt.Errorf("failed to parse user id: empty uid")
	}
	return uidStr, nil
}

func ParseUserUidStrFromTokenWithOldJwt(token *jwtOld.Token) (uid string, err error) {
	claims, ok := token.Claims.(jwtOld.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to convert token claims to jwt")
	}
	uidStr := fmt.Sprintf("%v", claims[JWTUserId])
	if uidStr == "" {
		return "", fmt.Errorf("failed to parse user id: empty uid")
	}
	return uidStr, nil
}
