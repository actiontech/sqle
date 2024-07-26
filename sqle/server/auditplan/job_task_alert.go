package auditplan

import (
	"time"

	"github.com/actiontech/sqle/sqle/server"
	"github.com/sirupsen/logrus"
)

type AuditPlanAlertJob struct {
	server.BaseJob
}

func NewAuditPlanAlertJob(entry *logrus.Entry) server.ServerJob {
	entry = entry.WithField("job", "audit_plan_alert")
	j := &AuditPlanAlertJob{}
	j.BaseJob = *server.NewBaseJob(entry, 5*time.Second, j.Alert)
	return j
}

func (j *AuditPlanAlertJob) Alert(entry *logrus.Entry) {
	return
}
