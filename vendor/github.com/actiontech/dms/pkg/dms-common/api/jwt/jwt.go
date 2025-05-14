package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	dmsCommonV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"

	"github.com/actiontech/dms/internal/dms/pkg/constant"
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
	JWTLoginType     = "loginType"
	JWTType          = "typ"

	DefaultDmsTokenExpHours        = 2
	DefaultDmsRefreshTokenExpHours = 24
)

func GenJwtToken(customClaims ...CustomClaimFunc) (tokenStr string, err error) {
	mapClaims := jwt.MapClaims{
		"iss":          "actiontech dms",
		JWTExpiredTime: jwt.NewNumericDate(time.Now().Add(DefaultDmsTokenExpHours * time.Hour)),
		JWTType:        constant.DMSToken,
	}

	return genJwtToken(mapClaims, customClaims...)
}

func GenJwtTokenWithExpirationTime(expiredTime *jwt.NumericDate, customClaims ...CustomClaimFunc) (tokenStr string, err error) {
	mapClaims := jwt.MapClaims{
		"iss":          "actiontech dms",
		JWTExpiredTime: expiredTime,
	}

	return genJwtToken(mapClaims, customClaims...)
}

func GenRefreshToken(customClaims ...CustomClaimFunc) (tokenStr string, err error) {
	mapClaims := jwt.MapClaims{
		"iss":          "actiontech dms",
		JWTExpiredTime: jwt.NewNumericDate(time.Now().Add(DefaultDmsRefreshTokenExpHours * time.Hour)),
		JWTType:        constant.DMSRefreshToken,
	}

	return genJwtToken(mapClaims, customClaims...)
}

func ParseRefreshToken(tokenStr string) (userUid, sub, sid string, expired bool, err error) {
	token, err := parseJwtTokenStr(tokenStr)
	if err != nil {
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) && validationErr.Errors&jwt.ValidationErrorExpired != 0 {
			expired = true
		} else {
			return "", "", "", false, err
		}
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", "", expired, fmt.Errorf("failed to convert token claims to jwt")
	}

	if fmt.Sprint(claims[JWTType]) != constant.DMSRefreshToken {
		return "", "", "", expired, fmt.Errorf("invalid jwt type")
	}

	userUid, _ = claims[JWTUserId].(string)
	if userUid == "" {
		return "", "", "", expired, fmt.Errorf("failed to parse user id: empty userUid")
	}
	sub, _ = claims["sub"].(string)
	sid, _ = claims["sid"].(string)

	return userUid, sub, sid, expired, nil
}

func genJwtToken(mapClaims jwt.MapClaims, customClaims ...CustomClaimFunc) (tokenStr string, err error) {
	for _, claimFunc := range customClaims {
		claimFunc(mapClaims)
	}
	mapClaims["iat"] = jwt.NewNumericDate(time.Now())
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

func WithJTWType(typ string) CustomClaimFunc {
	return func(claims jwt.MapClaims) {
		claims[JWTType] = typ
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

func WithAccessTokenMark(loginType string) CustomClaimFunc {
	return func(claims jwt.MapClaims) {
		claims[JWTLoginType] = loginType
	}
}

func WithSub(sub string) CustomClaimFunc {
	return func(claims jwt.MapClaims) {
		claims["sub"] = sub
	}
}

func WithSid(sid string) CustomClaimFunc {
	return func(claims jwt.MapClaims) {
		claims["sid"] = sid
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
		return token, fmt.Errorf("parse token failed: %w", err)
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

// 获取token的过期时间
func ParseExpiredTimeFromJwtTokenStr(tokenStr string) (expiredTime int64, err error) {
	// 使用自定义解析器，跳过过期验证
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, err := parser.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if signMethod256, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		} else if signMethod256 != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}

		return dmsCommonV1.JwtSigningKey, nil
	})
	if err != nil {
		return 0, fmt.Errorf("parse token failed: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("failed to convert token claims to jwt")
	}
	expiredTimeStr, ok := claims[JWTExpiredTime]
	if !ok {
		return 0, jwt.NewValidationError("unknown token", jwt.ValidationErrorClaimsInvalid)
	}
	expiredTimeInt, ok := expiredTimeStr.(float64)
	if !ok {
		return 0, jwt.NewValidationError("unknown token", jwt.ValidationErrorClaimsInvalid)
	}
	return int64(expiredTimeInt), nil
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

type TokenDetail struct {
	TokenStr  string
	UID       string
	LoginType string
}

// 由于sqle使用的github.com/golang-jwt/jwt，本方法为sqle兼容
func GetTokenDetailFromContextWithOldJwt(c EchoContextGetter) (tokenDetail *TokenDetail, err error) {
	tokenDetail = &TokenDetail{}

	if c.Get("user") == nil {
		return tokenDetail, nil
	}

	// Gets user token from the context.
	u, ok := c.Get("user").(*jwtOld.Token)
	if !ok {
		return nil, fmt.Errorf("failed to convert user from jwt token")
	}
	tokenDetail.TokenStr = u.Raw

	// get uid from token
	uid, err := ParseUserUidStrFromTokenWithOldJwt(u)
	if err != nil {
		return nil, err
	}
	tokenDetail.UID = uid

	// get login type from token
	claims, ok := u.Claims.(jwtOld.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to convert token claims to jwt")
	}
	loginType, ok := claims[JWTLoginType]
	if !ok {
		return tokenDetail, nil
	}

	tokenDetail.LoginType = fmt.Sprint(loginType)
	return tokenDetail, nil
}

func GetTokenDetailFromContext(c EchoContextGetter) (tokenDetail *TokenDetail, err error) {
	tokenDetail = &TokenDetail{}
	if c.Get("user") == nil {
		return tokenDetail, nil
	}

	// Gets user token from the context.
	u, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, fmt.Errorf("failed to convert user from jwt token")
	}
	tokenDetail.TokenStr = u.Raw

	// get uid from token
	uid, err := ParseUserUidStrFromToken(u)
	if err != nil {
		return nil, err
	}
	tokenDetail.UID = uid

	// get login type from token
	claims, ok := u.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to convert token claims to jwt")
	}
	loginType, ok := claims[JWTLoginType]
	if !ok {
		return tokenDetail, nil
	}

	tokenDetail.LoginType = fmt.Sprint(loginType)
	return tokenDetail, nil
}
