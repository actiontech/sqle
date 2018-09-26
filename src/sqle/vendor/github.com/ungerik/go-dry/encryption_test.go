package dry

import (
	"testing"
)

func Test_EncryptionAES(t *testing.T) {
	key := []byte("0123456789ABCDEF")
	data := "Hello World 1234"

	ciphertext := EncryptAES(key, []byte(data))
	if string(DecryptAES(key, ciphertext)) != data {
		t.Fail()
	}
}
