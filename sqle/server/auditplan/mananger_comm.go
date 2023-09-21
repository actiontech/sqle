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
}

func NewSyncFromAuditPlan(auditReport *model.AuditPlanReportV2, filterSqls []*model.AuditPlanSQLV2) SqlManager {
	return &SyncFromAuditPlan{
		AuditReport: auditReport,
		FilterSqls:  filterSqls,
	}
}

type SyncFromSqlAuditRecord struct {
	Task             *model.Task
	SqlFpMap         map[string]string
	ProjectId        uint
	SqlAuditRecordID uint
}

func NewSyncFromSqlAudit(task *model.Task, fpMap map[string]string, projectID uint, sqlAuditID uint) SqlManager {
	return &SyncFromSqlAuditRecord{
		Task:             task,
		ProjectId:        projectID,
		SqlAuditRecordID: sqlAuditID,
		SqlFpMap:         fpMap,
	}
}
