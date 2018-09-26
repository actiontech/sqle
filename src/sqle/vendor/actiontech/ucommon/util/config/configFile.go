package config

import (
	"actiontech/ucommon/conf"
	"actiontech/ucommon/log"
	"actiontech/ucommon/os"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	goOs "os"
	"strings"
)

func EncryptConfigPassword(raw string) string {
	return "encrypt_" + encrypt(raw)
}

func DecryptConfigPassword(encrypt string) string {
	if strings.HasPrefix(encrypt, "encrypt_") {
		return decrypt(strings.TrimPrefix(encrypt, "encrypt_"))
	}
	a, _ := base64.StdEncoding.DecodeString(encrypt)
	b, _ := base64.StdEncoding.DecodeString(string(a))
	return string(b)
}

func FixConfigFilePassword(stage *log.Stage, filePath string, config *conf.ConfigFile) error {
	sections := config.GetSections()
	for _, section := range sections {

		options, err := config.GetOptions(section)
		if nil != err {
			return err
		}
		for _, option := range options {
			if strings.HasSuffix(option, "password") {
				raw, _ := config.GetString(section, option)
				config.RemoveOption(section, option)
				config.AddOption(section, option+"_encrypt", EncryptConfigPassword(raw))

			} else if strings.HasSuffix(option, "password_encrypt") {
				raw, _ := config.GetString(section, option)
				if !strings.HasPrefix(raw, "encrypt_") {
					config.RemoveOption(section, option)
					config.AddOption(section, option+"_encrypt", EncryptConfigPassword(DecryptConfigPassword(raw)))
				}
			}
		}
	}

	if err := os.CopyFile(stage, filePath, filePath+".safe", "", 0750); nil != err {
		return err
	}

	if err := config.WriteConfigFile(filePath, 0640, "config password is encoded automatically, Version :2", []string{"global", "node1", "node2", "system"}); nil != err {
		goOs.Rename(filePath+".safe", filePath) //in case the disk is full
		return err
	} else {
		goOs.Remove(filePath + ".safe")
	}
	return nil
}

func DecryptConfigFilePassword(key, val string) string {
	if strings.HasSuffix(key, "_encrypt") {
		return DecryptConfigPassword(val)
	}
	return val
}

func encrypt(src string) string {
	block, iv := generateKey()
	encrypter := cipher.NewCFBEncrypter(block, iv)

	encrypted := make([]byte, len(src))
	encrypter.XORKeyStream(encrypted, []byte(src))
	return base64.StdEncoding.EncodeToString([]byte(encrypted))
}
func decrypt(src string) string {
	block, iv := generateKey()
	decrypter := cipher.NewCFBDecrypter(block, iv)

	if base64Deode, err := base64.StdEncoding.DecodeString(src); nil != err {
		panic(fmt.Sprintf("decrypt (%v) err (%v)", src, err))
	} else {
		src = string(base64Deode)
	}
	decrypted := make([]byte, len(src))
	decrypter.XORKeyStream(decrypted, []byte(src))
	return string(decrypted)
}
func generateKey() (cipher.Block, []byte) {
	block, err := aes.NewCipher([]byte("31940c058234fee9d825e8f9f0dc6e78"))
	if nil != err {
		panic(err.Error())
	}
	ciphertext := []byte("40be45e31b0d980a545f81ca7e336bc9a169f39b")
	iv := ciphertext[:aes.BlockSize] // const BlockSize = 16
	return block, iv
}
