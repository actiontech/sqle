package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"unsafe"
)

// base64 encoding string to decode string

func DecodeString(base64Str string) (string, error) {
	sDec, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}
	return *(*string)(unsafe.Pointer(&sDec)), nil
}

func Md5String(data string) string {
	md5 := md5.New()
	md5.Write([]byte(data))
	md5Data := md5.Sum([]byte(nil))
	return hex.EncodeToString(md5Data)
}
