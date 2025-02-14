//go:build enterprise
// +build enterprise

package license

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"
)

// TestGenerateLicense 测试生成许可证的函数
func TestGenerateKnowledgeBaseLicense(t *testing.T) {
	secretKey := []byte(secretKey)
	expirationDuration := 24 * time.Hour
	product := "TestProduct"

	license, err := generateLicense(expirationDuration, product, secretKey)
	if err != nil {
		t.Errorf("generateLicense returned an error: %v", err)
	}

	if license == "" {
		t.Error("generateLicense returned an empty license")
	}
}

// TestCheckLicenseValid 测试验证有效许可证的函数
func TestCheckKnowledgeBaseLicenseValid(t *testing.T) {
	secretKey := []byte(secretKey)
	expirationDuration := 24 * time.Hour * 10
	product := "TestProduct"

	license, err := generateLicense(expirationDuration, product, secretKey)
	if err != nil {
		t.Errorf("generateLicense returned an error: %v", err)
	}

	err = CheckKnowledgeBaseLicense(license)
	if err != nil {
		t.Error("CheckLicense should return true for a valid license")
	}

}

// TestCheckLicenseInvalidSignature 测试验证签名无效的许可证的函数
func TestCheckKnowledgeBaseLicenseInvalidSignature(t *testing.T) {
	secretKey := []byte(secretKey)
	expirationDuration := 24 * time.Hour
	product := "TestProduct"

	license, err := generateLicense(expirationDuration, product, secretKey)
	if err != nil {
		t.Errorf("generateLicense returned an error: %v", err)
	}

	// 篡改签名
	decodedLicense, err := base64.StdEncoding.DecodeString(license)
	if err != nil {
		t.Errorf("Base64 decode error: %v", err)
	}
	parts := strings.SplitN(string(decodedLicense), ".", 2)
	if len(parts) != 2 {
		t.Error("Invalid license format")
	}
	metaJSON := parts[0]
	fakeSignature := "fake_signature"
	fakeLicenseData := metaJSON + "." + fakeSignature
	fakeLicense := base64.StdEncoding.EncodeToString([]byte(fakeLicenseData))

	err = CheckKnowledgeBaseLicense(fakeLicense)
	if err != nil {
		t.Error("CheckLicense should return false for an invalid signature license")
	}
}

// TestCheckLicenseExpired 测试验证过期许可证的函数
func TestCheckKnowledgeBaseLicenseExpired(t *testing.T) {
	secretKey := []byte(secretKey)
	expirationDuration := -24 * time.Hour // 让许可证过期
	product := "TestProduct"

	license, err := generateLicense(expirationDuration, product, secretKey)
	if err != nil {
		t.Errorf("generateLicense returned an error: %v", err)
	}

	err = CheckKnowledgeBaseLicense(license)
	if err != nil {
		t.Error("CheckLicense should return false for an expired license")
	}
}
