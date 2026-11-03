package auditplan

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 轻量 TC3 签名客户端：仅覆盖本片所需 OpenAPI，避免全量 vendor 腾讯 SDK。
// 日志与错误包装禁止打印 SecretKey / Authorization / 签名串。

type tencentTC3Credential struct {
	SecretId  string
	SecretKey string
}

type tencentTC3Client struct {
	cred       tencentTC3Credential
	region     string
	service    string
	host       string
	version    string
	httpClient *http.Client
}

func newTencentTC3Client(secretId, secretKey, region, service, host, version string) *tencentTC3Client {
	return &tencentTC3Client{
		cred:       tencentTC3Credential{SecretId: secretId, SecretKey: secretKey},
		region:     region,
		service:    service,
		host:       host,
		version:    version,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *tencentTC3Client) doJSON(action string, payload interface{}, out interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}
	timestamp := time.Now().UTC().Unix()
	authorization := c.sign(action, body, timestamp)

	req, err := http.NewRequest(http.MethodPost, "https://"+c.host, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Host", c.host)
	req.Header.Set("X-TC-Action", action)
	req.Header.Set("X-TC-Timestamp", strconv.FormatInt(timestamp, 10))
	req.Header.Set("X-TC-Version", c.version)
	req.Header.Set("X-TC-Region", c.region)
	req.Header.Set("Authorization", authorization)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("tencent api request failed: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("tencent api http status %d", resp.StatusCode)
	}

	var envelope struct {
		Response json.RawMessage `json:"Response"`
	}
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return fmt.Errorf("decode response envelope: %w", err)
	}
	var apiErr struct {
		Error *struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error"`
		RequestId string `json:"RequestId"`
	}
	_ = json.Unmarshal(envelope.Response, &apiErr)
	if apiErr.Error != nil && apiErr.Error.Code != "" {
		// 不回传 Message 中可能含敏感信息的完整原文到上层日志以外；这里仅 code + request id
		return fmt.Errorf("tencent api error code=%s request_id=%s", apiErr.Error.Code, apiErr.RequestId)
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(envelope.Response, out); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}
	return nil
}

func (c *tencentTC3Client) sign(action string, payload []byte, timestamp int64) string {
	hashedPayload := sha256Hex(payload)
	canonicalHeaders := fmt.Sprintf("content-type:%s\nhost:%s\nx-tc-action:%s\n",
		"application/json; charset=utf-8", c.host, strings.ToLower(action))
	signedHeaders := "content-type;host;x-tc-action"
	canonicalRequest := strings.Join([]string{
		"POST",
		"/",
		"",
		canonicalHeaders,
		signedHeaders,
		hashedPayload,
	}, "\n")

	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, c.service)
	stringToSign := strings.Join([]string{
		"TC3-HMAC-SHA256",
		strconv.FormatInt(timestamp, 10),
		credentialScope,
		sha256Hex([]byte(canonicalRequest)),
	}, "\n")

	secretDate := hmacSHA256([]byte("TC3"+c.cred.SecretKey), date)
	secretService := hmacSHA256(secretDate, c.service)
	secretSigning := hmacSHA256(secretService, "tc3_request")
	signature := hex.EncodeToString(hmacSHA256(secretSigning, stringToSign))

	return fmt.Sprintf("TC3-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		c.cred.SecretId, credentialScope, signedHeaders, signature)
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	_, _ = h.Write([]byte(data))
	return h.Sum(nil)
}
