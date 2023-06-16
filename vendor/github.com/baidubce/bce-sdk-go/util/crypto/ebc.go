package crypto

import (
	"bytes"
	"crypto/aes"
	"fmt"
)

func EBCEncrypto(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	content := pKCS5Padding([]byte(data), block.BlockSize())
	crypted := make([]byte, len(content))

	if len(content)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("crypto/cipher: input not full blocks")
	}

	src := content
	dst := crypted
	for len(src) > 0 {
		block.Encrypt(dst, src[:block.BlockSize()])
		src = src[block.BlockSize():]
		dst = dst[block.BlockSize():]
	}

	return crypted, nil
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
