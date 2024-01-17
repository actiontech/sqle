package scanners

import (
	"context"
	"time"
)

type SQL struct {
	Fingerprint string
	RawText     string
	Counter     int
	Schema      string
	QueryTime   float64   // 慢日志执行时长
	QueryAt     time.Time // 慢日志发生时间
	DBUser      string    // 执行SQL的用户
	Endpoint    string    // 下发SQL的端点信息
}

// Scanner is a interface for all Scanners.
type Scanner interface {
	// Run start scanner. It parse SQLs and sends
	// it to the channel until ctx is canceled. Caller must check error.
	Run(ctx context.Context) error

	// SQLs returns channel that should be read until it is closed.
	SQLs() <-chan SQL

	// Upload upload sqls to underlying client.
	Upload(ctx context.Context, sqls []SQL) error
}
