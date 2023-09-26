package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
)

var SecretKey = []byte("471F77D078C5994BD06B65B8B5B1935B")

type encryptor struct {
	SecretKey []byte
	mutex     *sync.RWMutex
}

var std = &encryptor{
	SecretKey: SecretKey,
	mutex:     &sync.RWMutex{},
}

func NewEncryptor(key []byte) *encryptor {
	return &encryptor{
		SecretKey: key,
		mutex:     &sync.RWMutex{},
	}
}

func ResetAesSecretKey(secret string) error {
	if secret == "" {
		return nil
	}

	return std.SetAesSecretKey([]byte(secret))
}

func (e *encryptor) SetAesSecretKey(key []byte) (err error) {
	origKey := e.SecretKey
	e.SecretKey = key
	origData := "test"
	var secretData string
	defer func() {
		if err != nil {
			e.SecretKey = origKey
		}
	}()
	if secretData, err = e.AesEncrypt(origData); err != nil {
		return
	}
	if _, err = e.AesDecrypt(secretData); err != nil {
		return
	}
	return
}

func (e *encryptor) aesEncrypt(origData []byte) ([]byte, error) {
	e.mutex.RLock()
	key := e.SecretKey
	e.mutex.RUnlock()

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

func (e *encryptor) aesDecrypt(crypted []byte) ([]byte, error) {
	e.mutex.RLock()
	key := e.SecretKey
	e.mutex.RUnlock()

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

func (e *encryptor) AesEncrypt(origData string) (string, error) {
	crypted, err := e.aesEncrypt([]byte(origData))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func (e *encryptor) AesDecrypt(encrypted string) (string, error) {
	encryptedByte, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	origByte, err := e.aesDecrypt(encryptedByte)
	return string(origByte), err
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

func AesEncrypt(origData string) (string, error) {
	return std.AesEncrypt(origData)
}

func AesDecrypt(encrypted string) (string, error) {
	return std.AesDecrypt(encrypted)
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

func Md5(src string) string {
	m := md5.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}
