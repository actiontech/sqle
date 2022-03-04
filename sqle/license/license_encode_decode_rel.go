//go:build release
// +build release

package license

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	DELIMITER        = ";;"
	genEncoding      = base64.NewEncoding("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_~")
	collectEncoding  = base64.NewEncoding("012345ghijklmnopq6789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefrstuvwxyz_~")
	ExpireDateFormat = "2006-01-02"
)

const (
	// must be 16 byte
	cipherKey = "ActionTech--SQLE"

	// must be 16 byte
	cipherText = "a1b2c3d4e5f6g7h8"
)

var (
	ErrInvalidLicense     = errors.New("invalid license")
	ErrLicenseEmpty       = errors.New("license is empty")
	ErrCollectLicenseInfo = errors.New("collect license info error")
)

type LicensePermission struct {
	ExpireDate    string // License expire date eg: 2021-03-21
	Version       string // SQLE version
	UserCount     int    // The number of user
	InstanceCount int    // The number of instance
}

func NewLicensePermission(expireTime string, version string, userCount, instanceCount int) (*LicensePermission, error) {
	if _, err := time.Parse(ExpireDateFormat, expireTime); nil != err {
		return nil, fmt.Errorf("wrong expire date format(%v)", expireTime)
	}
	if userCount < 1 {
		return nil, errors.New("wrong user count")
	}
	if instanceCount < 1 {
		return nil, errors.New("wrong instance count")
	}
	return &LicensePermission{
		ExpireDate:    expireTime,
		Version:       version,
		UserCount:     userCount,
		InstanceCount: instanceCount,
	}, nil
}

func EncodeLicense(permission *LicensePermission, collectedInfosContent string) (string, error) {

	block, err := aes.NewCipher([]byte(cipherKey))
	if nil != err {
		return "", err
	}

	encrypter := cipher.NewCFBEncrypter(block, []byte(cipherText))

	permissionStr, err := json.Marshal(permission)
	if nil != err {
		return "", err
	}
	encodedPermissionStr, err := encode(string(permissionStr), encrypter)
	if nil != err {
		return "", err
	}

	encodedCollectedInfos, err := encode(collectedInfosContent, encrypter)
	if nil != err {
		return "", err
	}

	ret := make([]string, 0)
	licenseInfo := fmt.Sprintf("This license is for: %+v", permission)
	ret = append(ret, licenseInfo)
	ret = append(ret, encodedPermissionStr)
	ret = append(ret, encodedCollectedInfos)
	return strings.Join(ret, DELIMITER), nil
}

func encode(str string, encrypter cipher.Stream) (string, error) {
	encrypted := make([]byte, len(str))
	encrypter.XORKeyStream(encrypted, []byte(str))
	return genEncoding.EncodeToString(encrypted), nil
}

func DecodeLicense(license string) (permission *LicensePermission, collectedInfosContent string, err error) {
	block, err := aes.NewCipher([]byte(cipherKey))
	if nil != err {
		return nil, "", err
	}
	decrypter := cipher.NewCFBDecrypter(block, []byte(cipherText))

	options := strings.Split(license, DELIMITER)
	if len(options) < 3 { //licenseInfos;;permissions;;collectedInfos...
		return nil, "", ErrInvalidLicense
	}
	permissionStr, err := decode(options[1], decrypter)
	if nil != err {
		return nil, "", err
	}

	collectedInfosContent, err = decode(options[2], decrypter)
	if nil != err {
		return nil, "", err
	}

	permission = &LicensePermission{}
	err = json.Unmarshal([]byte(permissionStr), &permission)
	if nil != err {
		return nil, "", err
	}

	return permission, collectedInfosContent, nil
}

func decode(str string, decrypter cipher.Stream) (string, error) {
	a, err := genEncoding.DecodeString(str)
	if nil != err {
		return "", err
	}
	decrypted := make([]byte, len(a))
	decrypter.XORKeyStream(decrypted, a)
	return string(decrypted), nil
}
