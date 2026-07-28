package loginencryption

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"sync"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

const (
	ModeDisabled   = "disabled"
	ModeCompatible = "compatible"
	ModeRequired   = "required"

	AlgorithmSM2     = "SM2"
	CipherModeC1C3C2 = "C1C3C2"

	// maxCipherHexLen limits SM2 ciphertext hex length to avoid oversized payloads.
	maxCipherHexLen = 2048
)

// PublicInfo is returned by GET /v1/login/encryption.
type PublicInfo struct {
	Enable     bool   `json:"enable"`
	Algorithm  string `json:"algorithm"`
	CipherMode string `json:"cipher_mode"`
	PublicKey  string `json:"public_key"`
	KeyID      string `json:"key_id"`
}

// Manager holds login password encryption state.
type Manager struct {
	mode         string
	privateKey   *sm2.PrivateKey
	publicKeyHex string
	keyID        string
}

var (
	globalMu sync.RWMutex
	global   = &Manager{mode: ModeDisabled}
)

// Init loads SM2 private key according to mode.
// mode: disabled | compatible | required
// privateKeyPath is required when mode is not disabled.
func Init(mode, privateKeyPath string) error {
	mode = strings.TrimSpace(strings.ToLower(mode))
	if mode == "" {
		mode = ModeDisabled
	}

	switch mode {
	case ModeDisabled:
		setGlobal(&Manager{mode: ModeDisabled})
		return nil
	case ModeCompatible, ModeRequired:
		if strings.TrimSpace(privateKeyPath) == "" {
			return fmt.Errorf("login encryption private key path is required when mode is %s", mode)
		}
		priv, err := loadPrivateKey(privateKeyPath)
		if err != nil {
			return fmt.Errorf("load login encryption private key failed: %v", err)
		}
		pubHex := publicKeyToHex(&priv.PublicKey)
		m := &Manager{
			mode:         mode,
			privateKey:   priv,
			publicKeyHex: pubHex,
			keyID:        deriveKeyID(pubHex),
		}
		setGlobal(m)
		return nil
	default:
		return fmt.Errorf("invalid login encryption mode: %s", mode)
	}
}

// GetManager returns the current login encryption manager.
func GetManager() *Manager {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return global
}

func setGlobal(m *Manager) {
	globalMu.Lock()
	defer globalMu.Unlock()
	global = m
}

// PublicInfo returns public encryption metadata for clients.
func (m *Manager) PublicInfo() PublicInfo {
	if m == nil || m.mode == ModeDisabled || m.privateKey == nil {
		return PublicInfo{Enable: false}
	}
	return PublicInfo{
		Enable:     true,
		Algorithm:  AlgorithmSM2,
		CipherMode: CipherModeC1C3C2,
		PublicKey:  m.publicKeyHex,
		KeyID:      m.keyID,
	}
}

// Enabled reports whether login encryption is active.
func (m *Manager) Enabled() bool {
	return m != nil && m.mode != ModeDisabled && m.privateKey != nil
}

// Mode returns current encryption mode.
func (m *Manager) Mode() string {
	if m == nil {
		return ModeDisabled
	}
	return m.mode
}

// ResolvePassword returns plaintext password according to encryption mode.
// It never returns ciphertext or private key material in error messages.
func (m *Manager) ResolvePassword(plainPassword, encryptedPassword, keyID string) (string, error) {
	if m == nil || m.mode == ModeDisabled {
		if plainPassword == "" {
			return "", fmt.Errorf("password is required")
		}
		return plainPassword, nil
	}

	encryptedPassword = strings.TrimSpace(encryptedPassword)
	if encryptedPassword == "" {
		if m.mode == ModeRequired {
			return "", fmt.Errorf("encrypted password is required")
		}
		if plainPassword == "" {
			return "", fmt.Errorf("password is required")
		}
		return plainPassword, nil
	}

	if keyID != "" && keyID != m.keyID {
		return "", fmt.Errorf("invalid key_id")
	}
	if len(encryptedPassword) > maxCipherHexLen {
		return "", fmt.Errorf("encrypted password is too long")
	}

	plain, err := m.decryptPassword(encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("decrypt password failed")
	}
	if plain == "" {
		return "", fmt.Errorf("password is required")
	}
	return plain, nil
}

func (m *Manager) decryptPassword(cipherHex string) (string, error) {
	raw, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}
	if len(raw) == 0 {
		return "", fmt.Errorf("empty ciphertext")
	}
	// sm-crypto omits the leading 0x04; gmsm Decrypt(C1C3C2) expects it.
	if raw[0] != 0x04 {
		raw = append([]byte{0x04}, raw...)
	}
	// C1(64) + C3(32) + C2(>=1) + leading 0x04 => at least 98 bytes
	if len(raw) < 98 {
		return "", fmt.Errorf("ciphertext too short")
	}
	plain, err := sm2.Decrypt(m.privateKey, raw, sm2.C1C3C2)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func loadPrivateKey(path string) (*sm2.PrivateKey, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimSpace(string(content))
	if trimmed == "" {
		return nil, fmt.Errorf("private key file is empty")
	}

	if strings.Contains(trimmed, "BEGIN") {
		priv, err := x509.ReadPrivateKeyFromPem([]byte(trimmed), nil)
		if err != nil {
			return nil, err
		}
		return priv, nil
	}

	hexKey := strings.ReplaceAll(trimmed, "\n", "")
	hexKey = strings.ReplaceAll(hexKey, "\r", "")
	hexKey = strings.ReplaceAll(hexKey, " ", "")
	dBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("private key must be PEM or hex")
	}
	if len(dBytes) == 0 || len(dBytes) > 32 {
		return nil, fmt.Errorf("invalid private key length")
	}
	padded := make([]byte, 32)
	copy(padded[32-len(dBytes):], dBytes)

	curve := sm2.P256Sm2()
	priv := new(sm2.PrivateKey)
	priv.PublicKey.Curve = curve
	priv.D = new(big.Int).SetBytes(padded)
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(padded)
	if priv.PublicKey.X == nil || !curve.IsOnCurve(priv.PublicKey.X, priv.PublicKey.Y) {
		return nil, fmt.Errorf("invalid private key")
	}
	return priv, nil
}

func publicKeyToHex(pub *sm2.PublicKey) string {
	return hex.EncodeToString(elliptic.Marshal(pub.Curve, pub.X, pub.Y))
}

func deriveKeyID(publicKeyHex string) string {
	sum := sha256.Sum256([]byte(publicKeyHex))
	return hex.EncodeToString(sum[:8])
}

// EncryptForTest encrypts plaintext with current manager public key for unit tests.
func EncryptForTest(plain string) (cipherHex string, keyID string, err error) {
	m := GetManager()
	if !m.Enabled() {
		return "", "", fmt.Errorf("encryption is disabled")
	}
	cipher, err := sm2.Encrypt(&m.privateKey.PublicKey, []byte(plain), rand.Reader, sm2.C1C3C2)
	if err != nil {
		return "", "", err
	}
	return hex.EncodeToString(cipher), m.keyID, nil
}
