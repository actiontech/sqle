package scanners

import (
	"context"
)

type SQL struct {
	Fingerprint string
	RawText     string
	Counter     int
	Schema      string
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
