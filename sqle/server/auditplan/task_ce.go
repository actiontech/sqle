//go:build !enterprise
// +build !enterprise

package auditplan

import (
	"github.com/sirupsen/logrus"
)

func NewSlowLogTask(entry *logrus.Entry, ap *AuditPlan) Task {
	return NewDefaultTask(entry, ap)
}
