//go:build enterprise
// +build enterprise

package license

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

/*
	本文件是临时的，用于验证 license，在正式版本中该部分会被删除
	需要验证的 license 可以在 sqle的配置文件中配置
*/

// LicenseMeta 定义许可证的元数据
type LicenseMeta struct {
	Expiration int64  `json:"expiration"`
	Product    string `json:"product"`
	// 可以根据需要添加更多元数据字段
}

const secretKey string = "knowledge_base_secret_key_/github.com/actiontech/sqle/"

// 生成签名
func generateSignature(data []byte, secretKey []byte) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// 生成 license
func generateLicense(expirationDuration time.Duration, product string, secretKey []byte) (string, error) {
	// 计算过期时间的时间戳
	expirationTime := time.Now().Add(expirationDuration).Unix()

	// 构建许可证元数据
	meta := LicenseMeta{
		Expiration: expirationTime,
		Product:    product,
	}

	// 将元数据转换为 JSON 格式
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return "", err
	}

	// 生成签名
	signature := generateSignature(metaJSON, secretKey)

	// 将元数据和签名用点号连接，并进行 Base64 编码
	licenseData := string(metaJSON) + "." + signature
	encodedLicense := base64.StdEncoding.EncodeToString([]byte(licenseData))

	return encodedLicense, nil
}

var ErrKnowledgeBaseLicenseInvalid = fmt.Errorf("knowledge base license invalid")

// 检查 license 是否有效
func CheckKnowledgeBaseLicense(license string) error {
	// 对许可证进行 Base64 解码
	decodedLicense, err := base64.StdEncoding.DecodeString(license)
	if err != nil {
		return ErrKnowledgeBaseLicenseInvalid
	}

	// 分割元数据和签名
	parts := strings.SplitN(string(decodedLicense), ".", 2)
	if len(parts) != 2 {
		return ErrKnowledgeBaseLicenseInvalid
	}

	metaJSON := parts[0]
	signature := parts[1]

	// 验证签名
	expectedSignature := generateSignature([]byte(metaJSON), []byte(secretKey))
	if expectedSignature != signature {
		return ErrKnowledgeBaseLicenseInvalid
	}

	// 解析元数据
	var meta LicenseMeta
	err = json.Unmarshal([]byte(metaJSON), &meta)
	if err != nil {
		return ErrKnowledgeBaseLicenseInvalid
	}

	// 检查许可证是否过期
	if time.Now().Unix() > meta.Expiration {
		return ErrKnowledgeBaseLicenseInvalid
	}
	return nil
}
