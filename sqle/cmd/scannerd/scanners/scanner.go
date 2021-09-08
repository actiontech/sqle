package scanners

import (
	"context"

	"actiontech.cloud/sqle/sqle/sqle/model"
)

// Scanner is a interface for all Scanners.
type Scanner interface {
	// Run start scanner. It parse SQLs and sends
	// it to the channel until ctx is canceled.
	Run(ctx context.Context) error

	// SQLs returns channel that should be read until it is closed.
	SQLs() <-chan model.AuditPlanSQL
}
