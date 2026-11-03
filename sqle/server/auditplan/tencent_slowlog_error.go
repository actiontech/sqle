package auditplan

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// S3 AC-019 稳定错误码（产品可定位；与 S2 致命/可重试行为分层，本文件只负责码+脱敏文案）。
const (
	TencentErrCodeAuthInvalid      = "AUTH_INVALID"
	TencentErrCodeCamDenied        = "CAM_DENIED"
	TencentErrCodeRegionMismatch   = "REGION_MISMATCH"
	TencentErrCodeInstanceNotFound = "INSTANCE_NOT_FOUND"
	TencentErrCodeRateLimited      = "RATE_LIMITED"
	TencentErrCodeNetworkOrService = "NETWORK_OR_SERVICE"
	TencentErrCodeRetentionExceed  = "RETENTION_EXCEEDED"
	TencentErrCodePageOrParse      = "PAGE_OR_PARSE"
	TencentErrCodeConfigInvalid    = "CONFIG_INVALID"
)

// 用户可见短文案（固定口径；禁止拼接可能含敏感参数的云原始 Message）。
var tencentErrUserMessages = map[string]string{
	TencentErrCodeAuthInvalid:      "腾讯云鉴权失败，请检查 SecretId/SecretKey",
	TencentErrCodeCamDenied:        "腾讯云 CAM 权限不足，请授予慢日志查询权限",
	TencentErrCodeRegionMismatch:   "Region 与实例不匹配，请检查 Region 配置",
	TencentErrCodeInstanceNotFound: "腾讯云实例不存在或不可访问",
	TencentErrCodeRateLimited:      "腾讯云 API 限流，将在下一采集周期重试",
	TencentErrCodeNetworkOrService: "网络或腾讯云服务异常，将在下一采集周期重试",
	TencentErrCodeRetentionExceed:  "查询时间超出云侧慢日志保留期",
	TencentErrCodePageOrParse:      "分页或响应解析失败，请查看采集日志",
	TencentErrCodeConfigInvalid:    "扫描任务配置不完整",
}

// tencentCollectError 包装稳定码 + 脱敏短文案；Cause 可 Unwrap，供 S2 分类/测试匹配。
type tencentCollectError struct {
	Code    string
	Message string
	Cause   error
}

func (e *tencentCollectError) Error() string {
	if e == nil {
		return ""
	}
	msg := sanitizeTencentErrorText(e.Message)
	if e.Cause == nil {
		return fmt.Sprintf("[%s] %s", e.Code, msg)
	}
	// Cause 经 sanitize 后再拼，避免日志链路泄密；仍保留 unwrap。
	return fmt.Sprintf("[%s] %s (origin: %s)", e.Code, msg, sanitizeTencentErrorText(e.Cause.Error()))
}

func (e *tencentCollectError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func (e *tencentCollectError) ErrorCode() string {
	if e == nil {
		return ""
	}
	return e.Code
}

// tencentSQLMaxLen 单条 SQL 文本上限（对齐 hive slowlog 默认；超限优先跳过，AC-020）。
const tencentSQLMaxLen = 65535

// 脱敏：剥掉可能误入日志的密钥 / 签名片段。
var (
	reTencentSecretKey   = regexp.MustCompile(`(?i)(secret[_-]?key|access_key_secret)\s*[=:]\s*\S+`)
	reTencentAuthHeader  = regexp.MustCompile(`(?i)authorization\s*[:=]\s*\S+`)
	reTencentTC3Sig      = regexp.MustCompile(`(?i)TC3-HMAC-SHA256[^\s]*`)
	reTencentSignedChunk = regexp.MustCompile(`(?i)Signature\s*=\s*[A-Za-z0-9+/=]+`)
)

func sanitizeTencentErrorText(s string) string {
	if s == "" {
		return s
	}
	out := reTencentSecretKey.ReplaceAllString(s, "${1}=***")
	out = reTencentAuthHeader.ReplaceAllString(out, "Authorization=***")
	out = reTencentTC3Sig.ReplaceAllString(out, "TC3-HMAC-SHA256***")
	out = reTencentSignedChunk.ReplaceAllString(out, "Signature=***")
	return out
}

// classifyTencentCollectError 从 SDK/HTTP/本地错误映射稳定码 + 固定短文案。
// 分支顺序：配置 → 鉴权 → CAM → 实例不存在 → Region → 限流 → 保留期 → 分页/解析 → 默认网络/服务。
func classifyTencentCollectError(err error) (code, userMsg string) {
	if err == nil {
		return "", ""
	}
	s := err.Error()
	switch {
	case strings.Contains(s, "access_key_id is not configured"),
		strings.Contains(s, "access_key_secret is not configured"),
		strings.Contains(s, "region is not configured"),
		strings.Contains(s, "db_instance_id is not configured"):
		code = TencentErrCodeConfigInvalid
	case containsAnyFold(s,
		"AuthFailure",
		"SignatureFailure",
		"SecretIdNotFound",
		"CamSigAndAuthError",
		"UnauthorizedOperation",
	):
		code = TencentErrCodeAuthInvalid
	case containsAnyFold(s,
		"CamDeniedError",
		"InstanceAccessDeniedError",
		"OperationDenied.CamDenied",
	):
		code = TencentErrCodeCamDenied
	case containsAnyFold(s,
		"CdbInstanceNotFoundError",
		"InstanceNotFoundError",
		"ResourceNotFound.CdbInstanceNotFound",
		"ResourceNotFound.InstanceNotFound",
	):
		code = TencentErrCodeInstanceNotFound
	case containsAnyFold(s,
		"InvalidRegion",
		"RegionNotFound",
		"UnsupportedRegion",
		"InvalidParameterValue.Region",
		"InvalidParameter.Region",
		"region not match",
		"RegionMismatch",
	):
		code = TencentErrCodeRegionMismatch
	case containsAnyFold(s,
		"RequestLimitExceeded",
		"RequestLimit",
		"LimitExceeded",
	) && !isTencentResultOverLimitErr(err):
		code = TencentErrCodeRateLimited
	case containsAnyFold(s,
		"exceed retention",
		"Retention",
		"beyond the retention",
		"日志保留",
		"保留期",
	):
		code = TencentErrCodeRetentionExceed
	case errors.Is(err, errTencentMinWindowExceeded),
		isTencentResultOverLimitErr(err),
		containsAnyFold(s,
			"decode response",
			"marshal request",
			"分页",
			"min window still over limit",
		):
		code = TencentErrCodePageOrParse
	default:
		code = TencentErrCodeNetworkOrService
	}
	userMsg = tencentErrUserMessages[code]
	if userMsg == "" {
		userMsg = tencentErrUserMessages[TencentErrCodeNetworkOrService]
		code = TencentErrCodeNetworkOrService
	}
	return code, userMsg
}

func containsAnyFold(haystack string, needles ...string) bool {
	for _, n := range needles {
		if n != "" && strings.Contains(haystack, n) {
			return true
		}
	}
	return false
}

// wrapTencentCollectError 在 ExtractSQL 出口包装稳定码；保留 sentinel（禁用停请求）供 errors.Is。
func wrapTencentCollectError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, errTencentCollectStopped) {
		return err
	}
	var already *tencentCollectError
	if errors.As(err, &already) {
		return err
	}
	code, msg := classifyTencentCollectError(err)
	return &tencentCollectError{Code: code, Message: msg, Cause: err}
}

// tencentCollectErrorCode 提取稳定码（已包装或可分类）。
func tencentCollectErrorCode(err error) string {
	if err == nil {
		return ""
	}
	var te *tencentCollectError
	if errors.As(err, &te) {
		return te.Code
	}
	code, _ := classifyTencentCollectError(err)
	return code
}
