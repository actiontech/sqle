//go:build !enterprise
// +build !enterprise

package auditplan

import (
	"github.com/actiontech/sqle/sqle/model"

	"github.com/sirupsen/logrus"
)

func NewSlowLogTask(entry *logrus.Entry, ap *model.AuditPlan) Task {
	return NewDefaultTask(entry, ap)
}
