package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

var SecretKey = []byte("471F77D078C5994BD06B65B8B5B1935B")

func SetSecretKey(key []byte) (err error) {
	origKey := SecretKey
	SecretKey = key
	origData := "test"
	var secretData string
	defer func() {
		if err != nil {
			SecretKey = origKey
		}
	}()
	if secretData, err = AesEncrypt(origData); err != nil {
		return
	}
	if _, err = AesDecrypt(secretData); err != nil {
		return
	}
	return
}

func pKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	// if origData is empty, get byte by index length -1 will panic before go 11
	if length <= 0 {
		return origData
	}
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func aesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func aesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pKCS7UnPadding(origData)
	return origData, nil
}

func AesEncrypt(origData string) (string, error) {
	crypted, err := aesEncrypt([]byte(origData), SecretKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func AesDecrypt(encrypted string) (string, error) {
	encryptedByte, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	origByte, err := aesDecrypt(encryptedByte, SecretKey)
	return string(origByte), err
}

type Password string

func (p *Password) MarshalJSON() ([]byte, error) {
	if *p == "" {
		return json.Marshal([]byte(*p))
	}
	value, err := AesEncrypt(string(*p))
	if nil != err {
		return nil, fmt.Errorf("encrypt error: %v", err)
	}
	return json.Marshal(value)
}
