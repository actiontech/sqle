package loginencryption

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
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

	// DefaultPrivateKeyPath is used when login_encryption_private_key_path is unset.
	DefaultPrivateKeyPath = "./etc/login_sm2_private.pem"

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
// Empty mode is treated as required (zero-config default).
// Empty privateKeyPath uses DefaultPrivateKeyPath and auto-generates PEM 0600 if missing.
// An explicitly configured path that does not exist fails startup (no auto-generate).
func Init(mode, privateKeyPath string) error {
	mode = strings.TrimSpace(strings.ToLower(mode))
	if mode == "" {
		mode = ModeRequired
	}

	switch mode {
	case ModeDisabled:
		setGlobal(&Manager{mode: ModeDisabled})
		return nil
	case ModeCompatible, ModeRequired:
		path := strings.TrimSpace(privateKeyPath)
		allowGenerate := false
		if path == "" {
			path = DefaultPrivateKeyPath
			allowGenerate = true
		}
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				if !allowGenerate {
					// Explicit path: never auto-generate; fail startup (S2).
					return fmt.Errorf("login encryption private key path invalid or not readable: %s", path)
				}
				if err := generateAndSavePrivateKey(path); err != nil {
					return fmt.Errorf("generate login encryption private key failed: %v", err)
				}
			} else {
				return fmt.Errorf("stat login encryption private key failed: %v", err)
			}
		}
		priv, err := loadPrivateKey(path)
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

// ResolveUsername returns plaintext username according to encryption mode.
// Same key / algorithm / key_id rules as ResolvePassword.
// It never returns ciphertext, plaintext username, or private key material in error messages.
func (m *Manager) ResolveUsername(plainUsername, encryptedUsername, keyID string) (string, error) {
	if m == nil || m.mode == ModeDisabled {
		if plainUsername == "" {
			return "", fmt.Errorf("username is required")
		}
		return plainUsername, nil
	}

	if m.mode == ModeRequired && strings.TrimSpace(plainUsername) != "" {
		return "", fmt.Errorf("plaintext username is not allowed")
	}

	encryptedUsername = strings.TrimSpace(encryptedUsername)
	if encryptedUsername == "" {
		if m.mode == ModeRequired {
			return "", fmt.Errorf("encrypted username is required")
		}
		if plainUsername == "" {
			return "", fmt.Errorf("username is required")
		}
		return plainUsername, nil
	}

	if keyID != "" && keyID != m.keyID {
		return "", fmt.Errorf("invalid key_id")
	}
	if len(encryptedUsername) > maxCipherHexLen {
		return "", fmt.Errorf("encrypted username is too long")
	}

	plain, err := m.decryptCipher(encryptedUsername)
	if err != nil {
		return "", fmt.Errorf("decrypt username failed")
	}
	if plain == "" {
		return "", fmt.Errorf("username is required")
	}
	return plain, nil
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

	if m.mode == ModeRequired && strings.TrimSpace(plainPassword) != "" {
		return "", fmt.Errorf("plaintext password is not allowed")
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

	plain, err := m.decryptCipher(encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("decrypt password failed")
	}
	if plain == "" {
		return "", fmt.Errorf("password is required")
	}
	return plain, nil
}

func (m *Manager) decryptCipher(cipherHex string) (string, error) {
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

func generateAndSavePrivateKey(path string) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}
	priv, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	pemBytes, err := x509.WritePrivateKeyToPem(priv, nil)
	if err != nil {
		return err
	}
	// O_EXCL: never overwrite an existing key file.
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(pemBytes); err != nil {
		return err
	}
	return f.Sync()
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
