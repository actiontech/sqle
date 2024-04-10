package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"

	"github.com/xen0n/go-workwx/internal/lowlevel/pkcs7"
)

type WorkwxPayload struct {
	Msg       []byte
	ReceiveID []byte
}

type WorkwxEncryptor struct {
	aesKey        []byte
	entropySource io.Reader
}

type WorkwxEncryptorOption interface {
	applyTo(x *WorkwxEncryptor)
}

type customEntropySource struct {
	inner io.Reader
}

func WithEntropySource(e io.Reader) WorkwxEncryptorOption {
	return &customEntropySource{inner: e}
}

func (o *customEntropySource) applyTo(x *WorkwxEncryptor) {
	x.entropySource = o.inner
}

var errMalformedEncodingAESKey = errors.New("malformed EncodingAESKey")

func NewWorkwxEncryptor(
	encodingAESKey string,
	opts ...WorkwxEncryptorOption,
) (*WorkwxEncryptor, error) {
	aesKey, err := base64.StdEncoding.DecodeString(encodingAESKey + "=")
	if err != nil {
		return nil, err
	}

	if len(aesKey) != 32 {
		return nil, errMalformedEncodingAESKey
	}

	obj := WorkwxEncryptor{
		aesKey:        aesKey,
		entropySource: rand.Reader,
	}
	for _, o := range opts {
		o.applyTo(&obj)
	}

	return &obj, nil
}

func (e *WorkwxEncryptor) Decrypt(base64Msg []byte) (WorkwxPayload, error) {
	// base64 decode
	buflen := base64.StdEncoding.DecodedLen(len(base64Msg))
	buf := make([]byte, buflen)
	n, err := base64.StdEncoding.Decode(buf, base64Msg)
	if err != nil {
		return WorkwxPayload{}, err
	}
	buf = buf[:n]

	// init cipher
	block, err := aes.NewCipher(e.aesKey)
	if err != nil {
		return WorkwxPayload{}, err
	}

	iv := e.aesKey[:16]
	state := cipher.NewCBCDecrypter(block, iv)

	// decrypt in-place in the allocated temp buffer
	state.CryptBlocks(buf, buf)
	buf = pkcs7.Unpad(buf)

	// assemble decrypted payload
	// drop the 16-byte random prefix
	msglen := binary.BigEndian.Uint32(buf[16:20])
	msg := buf[20 : 20+msglen]
	receiveID := buf[20+msglen:]

	return WorkwxPayload{
		Msg:       msg,
		ReceiveID: receiveID,
	}, nil
}

func (e *WorkwxEncryptor) prepareBufForEncryption(payload *WorkwxPayload) ([]byte, error) {
	resultMsgLen := 16 + 4 + len(payload.Msg) + len(payload.ReceiveID)

	// allocate buffer
	buf := make([]byte, 16, resultMsgLen)

	// add random prefix
	_, err := io.ReadFull(e.entropySource, buf) // len(buf) == 16 at this moment
	if err != nil {
		return nil, err
	}

	buf = buf[:cap(buf)] // grow to full capacity
	binary.BigEndian.PutUint32(buf[16:], uint32(len(payload.Msg)))
	copy(buf[20:], payload.Msg)
	copy(buf[20+len(payload.Msg):], payload.ReceiveID)

	return pkcs7.Pad(buf), nil
}

func (e *WorkwxEncryptor) Encrypt(payload *WorkwxPayload) (string, error) {
	buf, err := e.prepareBufForEncryption(payload)
	if err != nil {
		return "", err
	}

	// init cipher
	block, err := aes.NewCipher(e.aesKey)
	if err != nil {
		return "", err
	}

	iv := e.aesKey[:16]
	state := cipher.NewCBCEncrypter(block, iv)

	// encrypt in-place as we own the buffer
	state.CryptBlocks(buf, buf)

	return base64.StdEncoding.EncodeToString(buf), nil
}
