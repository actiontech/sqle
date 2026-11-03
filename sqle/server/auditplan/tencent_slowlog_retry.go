package auditplan

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

const (
	tencentAPIMaxAttempts   = 5
	tencentAPIBackoffBase   = 200 * time.Millisecond
	tencentAPIBackoffCap    = 2 * time.Second
)

// isTencentFatalAPIErr 鉴权/CAM/实例不存在/参数错误 → 立即失败，不盲重试。
func isTencentFatalAPIErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	codes := []string{
		"AuthFailure",
		"UnauthorizedOperation",
		"SecretIdNotFound",
		"SignatureFailure",
		"CamSigAndAuthError",
		"CamDeniedError",
		"OperationDenied.CamDeniedError",
		"CdbInstanceNotFoundError",
		"ResourceNotFound.CdbInstanceNotFoundError",
		"InstanceNotFoundError",
		"ResourceNotFound.InstanceNotFoundError",
		"InstanceAccessDeniedError",
		"OperationDenied.InstanceAccessDeniedError",
		"OperationDenied.NotSupportBasic",
		"InvalidParameter",
		"InvalidParameterValue",
		"MissingParameter",
		"UnknownParameter",
	}
	for _, c := range codes {
		if strings.Contains(msg, c) {
			return true
		}
	}
	return false
}

// isTencentRetryableAPIErr 限流 / TimeoutError / 5xx / 网络超时等可重试。
func isTencentRetryableAPIErr(err error) bool {
	if err == nil {
		return false
	}
	if isTencentResultOverLimitErr(err) || isTencentFatalAPIErr(err) {
		return false
	}
	msg := err.Error()
	retrySignals := []string{
		"RequestLimitExceeded",
		"RequestLimit",
		"LimitExceeded",
		"FailedOperation.TimeoutError",
		"TimeoutError",
		"tencent api http status 5",
		"tencent api request failed",
		"i/o timeout",
		"connection reset",
		"connection refused",
		"TLS handshake timeout",
		"temporary failure",
		"EOF",
	}
	for _, s := range retrySignals {
		if strings.Contains(msg, s) {
			return true
		}
	}
	var netErr net.Error
	if errorsAsNet(err, &netErr) && netErr.Timeout() {
		return true
	}
	return false
}

func errorsAsNet(err error, target *net.Error) bool {
	for err != nil {
		if ne, ok := err.(net.Error); ok {
			*target = ne
			return true
		}
		type unwrapper interface{ Unwrap() error }
		u, ok := err.(unwrapper)
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
}

func tencentAPIBackoffDelay(attempt int, base, capDelay time.Duration, rnd *rand.Rand) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	d := base << attempt // 200ms, 400ms, 800ms, ...
	if d > capDelay || d <= 0 {
		d = capDelay
	}
	// 抖动：0.5x～1.0x
	if rnd == nil {
		rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	factor := 0.5 + rnd.Float64()*0.5
	return time.Duration(float64(d) * factor)
}

// withTencentAPIRetry 对单次 OpenAPI 调用做有限退避；缩窗类/致命类不进入退避。
// shouldStop 非 nil 时，每次 attempt 与 sleep 前检查；为 true 则停止新请求（AC-015），不进入收尾重试。
func withTencentAPIRetry(sleep func(time.Duration), maxAttempts int, shouldStop func() bool, fn func() error) error {
	if maxAttempts <= 0 {
		maxAttempts = tencentAPIMaxAttempts
	}
	if sleep == nil {
		sleep = time.Sleep
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if shouldStop != nil && shouldStop() {
			return errTencentCollectStopped
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if isTencentResultOverLimitErr(lastErr) || isTencentFatalAPIErr(lastErr) {
			return lastErr
		}
		if !isTencentRetryableAPIErr(lastErr) {
			return lastErr
		}
		if attempt+1 >= maxAttempts {
			break
		}
		if shouldStop != nil && shouldStop() {
			return errTencentCollectStopped
		}
		sleep(tencentAPIBackoffDelay(attempt, tencentAPIBackoffBase, tencentAPIBackoffCap, rnd))
	}
	return fmt.Errorf("tencent api retry exhausted after %d attempts: %w", maxAttempts, lastErr)
}
