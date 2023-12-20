package auditplan

import (
	"github.com/actiontech/sqle/sqle/model"
)

type SqlManager interface {
	SyncSqlManager() error
}

type SyncFromAuditPlan struct {
	AuditReport *model.AuditPlanReportV2
	FilterSqls  []*model.AuditPlanSQLV2
	Task        *model.Task
}

func NewSyncFromAuditPlan(auditReport *model.AuditPlanReportV2, filterSqls []*model.AuditPlanSQLV2, task *model.Task) SqlManager {
	return &SyncFromAuditPlan{
		AuditReport: auditReport,
		FilterSqls:  filterSqls,
		Task:        task,
	}
}

type SyncFromSqlAuditRecord struct {
	Task             *model.Task
	SqlFpMap         map[string]string
	ProjectId        string
	SqlAuditRecordID uint
}

func NewSyncFromSqlAudit(task *model.Task, fpMap map[string]string, projectID string, sqlAuditID uint) SqlManager {
	return &SyncFromSqlAuditRecord{
		Task:             task,
		ProjectId:        projectID,
		SqlAuditRecordID: sqlAuditID,
		SqlFpMap:         fpMap,
	}
}
